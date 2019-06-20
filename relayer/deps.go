package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// use these standard error codes to send out to request replies so callers can figure
// out what's going wrong with their requests
const (
	StatusOK                  = http.StatusOK                  //  200 RFC 7231, 6.3.1
	StatusBadRequest          = http.StatusBadRequest          //  400 RFC 7231, 6.5.1
	StatusUnauthorized        = http.StatusUnauthorized        //  401 RFC 7235, 3.1
	StatusNotFound            = http.StatusNotFound            //  404 RFC 7231, 6.5.4
	StatusInternalServerError = http.StatusInternalServerError //  500 RFC 7231, 6.6.1
	StatusBadGateway          = http.StatusBadGateway          //  502 RFC 7231, 6.6.3
)

// StatusResponse defines a generic status response structure
type StatusResponse struct {
	Code   int
	Status string
}

// checkGet checks if the invoming request is a GET request
func checkGet(w http.ResponseWriter, r *http.Request) {
	checkOrigin(w, r)
	if r.Method != "GET" {
		responseHandler(w, StatusNotFound)
		return
	}
}

// checkOrigin checks the origin of the incoming request
func checkOrigin(w http.ResponseWriter, r *http.Request) {
	// re-enable this function for all private routes
	// if r.Header.Get("Origin") != "localhost" { // allow only our frontend UI to connect to our RPC instance
	// 	http.Error(w, "404 page not found", http.StatusNotFound)
	// }
}

// responseHandler is the default response handler that sends out response codes on successful
// completion of certain calls
func responseHandler(w http.ResponseWriter, status int) {
	var response StatusResponse
	response.Code = status
	switch status {
	case StatusOK:
		response.Status = "OK"
	case StatusBadRequest:
		response.Status = "Bad Request error!"
	case StatusUnauthorized:
		response.Status = "You are unauthorized to make this request"
	case StatusNotFound:
		response.Status = "404 Error Not Found!"
	case StatusInternalServerError:
		response.Status = "Internal Server Error"
	case StatusBadGateway:
		response.Status = "Bad Gateway Error"
	default:
		response.Status = "404 Page Not Found"
	}
	MarshalSend(w, response)
}

// MarshalSend marshals and writes a json string into the writer
func MarshalSend(w http.ResponseWriter, x interface{}) {
	xJson, err := json.Marshal(x)
	if err != nil {
		log.Println("did not marshal json", err)
		errString := "Internal Server Error"
		WriteToHandler(w, []byte(errString))
		return
	}
	log.Println("JSON: ", string(xJson))
	WriteToHandler(w, xJson)
}

// WriteToHandler constructs a reply to the passed writer
func WriteToHandler(w http.ResponseWriter, jsonString []byte) {
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Authorization, Cache-Control, Content-Type")
	w.Header().Add("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}
