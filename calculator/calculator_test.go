package calculator

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidiateInput(t *testing.T) {
	t.Run("testing for int", func(t *testing.T) {
		postBody := []byte(`{"number1":1, "number2":2}`)
		req := httptest.NewRequest("POST", "/add/", bytes.NewBuffer(postBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ValidateInput(w, req)
		statusCheck(t, *w)
	})
	t.Run("testing for string", func(t *testing.T) {
		postBody := []byte(`{"number1":"1", "number2":"2"}`)
		req := httptest.NewRequest("POST", "/add/", bytes.NewBuffer(postBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		expectedError := "cannot unmarshal string"
		ValidateInput(w, req)
		assertError(t, *w, expectedError)
	})
	t.Run("testing for float", func(t *testing.T) {
		postBody := []byte(`{"number1":1.0, "number2":2.0}`)
		req := httptest.NewRequest("POST", "/add/", bytes.NewBuffer(postBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		expectedError := "cannot unmarshal number 1.0"
		ValidateInput(w, req)
		assertError(t, *w, expectedError)
	})
	t.Run("testing for incomplete num2", func(t *testing.T) {
		postBody := []byte(`{"number1":1}`)
		req := httptest.NewRequest("POST", "/add/", bytes.NewBuffer(postBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		expectedError := "Error:Field validation for 'Number2' failed"
		ValidateInput(w, req)
		assertError(t, *w, expectedError)
	})
	t.Run("testing for incomplete num1", func(t *testing.T) {
		postBody := []byte(`{"number2":1}`)
		req := httptest.NewRequest("POST", "/add/", bytes.NewBuffer(postBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		expectedError := "Error:Field validation for 'Number1' failed"
		ValidateInput(w, req)
		assertError(t, *w, expectedError)
	})
	t.Run("testing for incomplete num1 & num2", func(t *testing.T) {
		postBody := []byte(`{}`)
		req := httptest.NewRequest("POST", "/add/", bytes.NewBuffer(postBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		expectedError := "Error:Field validation for 'Number1' failed"
		ValidateInput(w, req)
		assertError(t, *w, expectedError)
	})
	t.Run("testing for zero division", func(t *testing.T) {
		postBody := []byte(`{"number1":1, "number2":0}`)
		req := httptest.NewRequest("POST", "/divide/", bytes.NewBuffer(postBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		expectedError := "Division by 0"
		ValidateInput(w, req)
		assertError(t, *w, expectedError)
	})
	t.Run("test add", func(t *testing.T) {
		postBody := []byte(`{"number1":1, "number2":2}`)
		req := httptest.NewRequest("POST", "/add/", bytes.NewBuffer(postBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ValidateInput(w, req)
		statusCheck(t, *w)
		var result Result
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil {
			t.Fatal(err)
		}
		assertValues(t, 3, result.Output)
	})
	t.Run("test subtract", func(t *testing.T) {
		postBody := []byte(`{"number1":2, "number2":1}`)
		req := httptest.NewRequest("POST", "/subtract/", bytes.NewBuffer(postBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ValidateInput(w, req)
		statusCheck(t, *w)
		var result Result
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil {
			t.Fatal(err)
		}
		assertValues(t, 1, result.Output)
	})
	t.Run("test multiply", func(t *testing.T) {
		postBody := []byte(`{"number1":2, "number2":3}`)
		req := httptest.NewRequest("POST", "/multiply/", bytes.NewBuffer(postBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ValidateInput(w, req)
		statusCheck(t, *w)
		var result Result
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil {
			t.Fatal(err)
		}
		assertValues(t, 6, result.Output)
	})
	t.Run("test divide", func(t *testing.T) {
		postBody := []byte(`{"number1":6, "number2":3}`)
		req := httptest.NewRequest("POST", "/divide/", bytes.NewBuffer(postBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ValidateInput(w, req)
		statusCheck(t, *w)
		var result Result
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil {
			t.Fatal(err)
		}
		assertValues(t, 2, result.Output)
	})
	t.Run("test torig", func(t *testing.T) {
		postBody := []byte(`{"number1":6, "number2":3}`)
		req := httptest.NewRequest("POST", "/troig/", bytes.NewBuffer(postBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ValidateInput(w, req)
		assertError(t, *w, "invalid path")
	})
}

func assertValues(t *testing.T, expected, got int) {
	t.Helper()
	if expected != got {
		t.Errorf("Expected %v got %v", expected, got)
	}
}

func assertError(t *testing.T, w httptest.ResponseRecorder, expectedError string) {
	t.Helper()
	if w.Code < 400 {
		t.Fatal("Expected error, did not get one.")
	}
	if !strings.Contains(w.Body.String(), expectedError) {
		t.Errorf("Expected: %v got: $%v", expectedError, w.Body.String())
	}
}

func statusCheck(t *testing.T, w httptest.ResponseRecorder) {
	t.Helper()
	if w.Code >= 400 {
		t.Fatalf("error: %v", w.Body)
	}
}
