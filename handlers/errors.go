package handlers

import (
	"errors"
	"net/http"
)

// validationErr is an error which occurred while validating user input,
// e.g. field cannot be less than x characters.
type validationErr struct {
	error
}

// isValidationErr checks to see if any error in the chain is a validationErr.
func isValidationErr(err error) bool {
	var vErr *validationErr
	return errors.As(err, &vErr)
}

// clientErr is an error which occurred while validating client input,
// e.g. request is missing a hidden form field.
type clientErr struct {
	error
}

func (c clientErr) Code() int {
	return http.StatusBadRequest
}
