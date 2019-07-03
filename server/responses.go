package server

import (
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
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
	erpc.MarshalSend(w, response)
}
