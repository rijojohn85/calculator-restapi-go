package main

import (
	"calculator/calculator"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// define middleware
type authenticationMiddleware struct {
	tokenUsers map[string]string
}

var timeDB = make(map[string]time.Time)

// initiliaze it
func (amw *authenticationMiddleware) Populate() {
	amw.tokenUsers["00000000"] = "user0"
	amw.tokenUsers["aaaaaaaa"] = "userA"
	amw.tokenUsers["05f717e5"] = "randomUser"
	amw.tokenUsers["deadbeef"] = "user0"
}

// Middleware function
func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")
		if user, found := amw.tokenUsers[token]; found {
			// log user to log file
			r.Header.Add("user", user)
			fmt.Printf("Authenticated User: %s\n", user)
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := httptest.NewRecorder()
		start := time.Now()
		user := r.Header.Get("user")
		handler.ServeHTTP(rec, r)
		w.WriteHeader(rec.Result().StatusCode)
		io.Copy(w, rec.Result().Body)
		red := rec.Body
		log.Printf("[Address: %s] [Request Method: %s] [Request URL: %s] [ExecutionTime: %v] [ResponseCode: %v] [user: %s] [ResponseBody: %v]\n", r.RemoteAddr, r.Method, r.URL, time.Since(start), rec.Result().StatusCode, user, red)
	})
}

func rateLimiter(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("user")
		if userLastLoginTime, found := timeDB[user]; found {
			if time.Since(userLastLoginTime) < time.Second*5 {
				http.Error(w, "please give 5 seconds between Request", http.StatusTooEarly)
			} else {
				timeDB[user] = time.Now()
				handler.ServeHTTP(w, r)
			}
		} else {
			timeDB[user] = time.Now()
			handler.ServeHTTP(w, r)
		}
	})
}

func main() {
	logPath := "development.log"
	httpPort := 8080
	openLogFile(logPath)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	router := mux.NewRouter()

	router.HandleFunc("/add/", calculator.ValidateInput).Methods("POST")
	router.HandleFunc("/subtract/", calculator.ValidateInput).Methods("POST")
	router.HandleFunc("/multiply/", calculator.ValidateInput).Methods("POST")
	router.HandleFunc("/divide/", calculator.ValidateInput).Methods("POST")
	amw := authenticationMiddleware{tokenUsers: make(map[string]string)}
	amw.Populate()

	var middleware []mux.MiddlewareFunc
	middleware = []mux.MiddlewareFunc{amw.Middleware, rateLimiter, logRequest}
	for _, m := range middleware {
		router.Use(m)
	}
	fmt.Printf("Listening on %v\n", httpPort)
	fmt.Printf("Logging to %v\n", logPath)
	err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), router)
	if err != nil {
		log.Fatal(err)
	}
}

// helper functions

/*
* opens file for logging and passes it to log
 */
func openLogFile(logfile string) {
	if logfile != "" {
		// open file are Read/Write, Create, append with permission 0640
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
		if err != nil {
			log.Fatal("Unable to openLogFile: ", err)
		}
		log.SetOutput(lf)
	}
}
