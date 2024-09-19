package calculator

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Number interface {
	int64 | float64
}

type InputParameters[T Number] struct {
	Number1 *T `json:"number1,omitempty" validate:"required"`
	Number2 *T `json:"number2,omitempty" validate:"required"`
}

type Result[T Number] struct {
	Output T `json:"output"`
}

func ValidateInput(w http.ResponseWriter, r *http.Request) {
	input := InputParameters{Number1: nil, Number2: nil}
	result := Result{Output: 0}
	validator := validator.New(validator.WithRequiredStructEnabled())

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	validationError := validator.Struct(input)
	if validationError != nil {
		http.Error(w, validationError.Error(), http.StatusBadRequest)
		return
	}
	path := r.URL.String()
	if path == "/divide/" && *input.Number2 == 0 {
		http.Error(w, "Division by 0", http.StatusBadRequest)
		return
	}
	// fmt.Printf("Num1: %d, Num2: %d, URL: %s\n", *input.Number1, *input.Number2, path)
	switch path {
	case "/add/":
		result.Output = *input.Number1 + *input.Number2
	case "/subtract/":
		result.Output = *input.Number1 - *input.Number2
	case "/multiply/":
		result.Output = *input.Number1 * *input.Number2
	case "/divide/":
		result.Output = *input.Number1 / *input.Number2
	default:
		{
			http.Error(w, "invalid path", http.StatusForbidden)
			return
		}
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(result)
}
