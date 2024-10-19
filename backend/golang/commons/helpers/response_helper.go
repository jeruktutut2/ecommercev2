package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type Response struct {
	Data   interface{} `json:"data"`
	Errors interface{} `json:"errors"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

type ErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ToErrorMessages(message string) (errorMessages []ErrorMessage) {
	var errorMessage ErrorMessage
	errorMessage.Field = "message"
	errorMessage.Message = message
	errorMessages = append(errorMessages, errorMessage)
	return
}

func toResponseRequestTimeout() (httpCode int, response Response) {
	message := "time out or user cancel the request"
	errorMessages := ToErrorMessages(message)
	httpCode = http.StatusRequestTimeout
	response = Response{
		Data:   nil,
		Errors: errorMessages,
	}
	return
}

func ToResponseInternalServerError() (httpCode int, response Response) {
	message := "internal server error"
	errorMessages := ToErrorMessages(message)
	httpCode = http.StatusInternalServerError
	response = Response{
		Data:   nil,
		Errors: errorMessages,
	}
	return
}

func ToResponseCheckError(err error, requestId string) (httpCode int, response Response) {
	PrintLogToTerminal(err, requestId)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return toResponseRequestTimeout()
	} else {
		return ToResponseInternalServerError()
	}
}

func ToResponseRequestValidation(requestId string, validationErrorMessages []ErrorMessage) (httpCode int, response Response) {
	var validationErrorMessageByte []byte
	validationErrorMessageByte, err := json.Marshal(validationErrorMessages)
	if err != nil {
		return ToResponseCheckError(err, requestId)
	}
	err = errors.New(string(validationErrorMessageByte))
	PrintLogToTerminal(err, requestId)
	httpCode = http.StatusBadRequest
	response = Response{
		Data:   nil,
		Errors: validationErrorMessages,
	}
	return
}

func ToResponseError(err error, requestId string, httpCode int, message string) (int, Response) {
	PrintLogToTerminal(err, requestId)
	errorMessages := ToErrorMessages(message)
	response := Response{
		Data:   nil,
		Errors: errorMessages,
	}
	return httpCode, response
}
