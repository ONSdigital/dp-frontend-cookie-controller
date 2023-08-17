package handlers

import (
	"net/http"
)

// clientErr is an error which occurred while validating client input,
// e.g. request is missing a hidden form field.
type clientErr struct {
	error
}

func (c clientErr) Code() int {
	return http.StatusBadRequest
}
