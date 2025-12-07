package main

import (
	"encoding/json"
	"net/http"
)

type calcRequest struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

type calcResult struct {
	Result float64 `json:"result"`
}

func validateRequest(req *http.Request) (float64, float64, error) {
	var reqBody calcRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		return -1, -1, err
	}

	return reqBody.A, reqBody.B, nil
}

func sendResponse(w http.ResponseWriter, result float64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := calcResult{
		Result: result,
	}

	json.NewEncoder(w).Encode(response)
}

func add(w http.ResponseWriter, req *http.Request) {
	a, b, err := validateRequest(req)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	result := a + b

	sendResponse(w, result)
}

func sub(w http.ResponseWriter, req *http.Request) {
	a, b, err := validateRequest(req)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	result := a - b

	sendResponse(w, result)
}

func multi(w http.ResponseWriter, req *http.Request) {
	a, b, err := validateRequest(req)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	result := a * b

	sendResponse(w, result)
}

func div(w http.ResponseWriter, req *http.Request) {
	a, b, err := validateRequest(req)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	result := a / b

	sendResponse(w, result)
}

func main() {
	http.HandleFunc("/add", add)
	http.HandleFunc("/sub", sub)
	http.HandleFunc("/multi", multi)
	http.HandleFunc("/div", div)

	http.ListenAndServe(":8080", nil)
}
