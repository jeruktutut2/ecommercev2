package middlewares

import (
	"backend-golang/commons/helpers"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type responseBodyWriter struct {
	io.Writer
	http.ResponseWriter
	status int
}

func (w *responseBodyWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func PrintRequestResponseLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		datetimeNowRequest := time.Now()
		requestMethod := c.Request().Method
		var err error
		var errorMessages []helpers.ErrorMessage
		var httpResponseCode int

		// why put this, because i wanna log request when there is no request id
		// requestId := c.Request().Header.Get("X-REQUEST-ID")
		// if requestId == "" {
		// 	errRequestId := errors.New("cannot find requestId")
		// 	helpers.PrintLogToTerminal(errRequestId, "requestId")
		// 	httpResponseCode = http.StatusBadRequest
		// 	errorMessages = helpers.ToErrorMessages("cannot find requestId")
		// 	err = errRequestId
		// }
		// ctx := context.WithValue(c.Request().Context(), RequestIdKey, requestId)
		// c.SetRequest(c.Request().WithContext(ctx))
		requestId := c.Request().Context().Value(RequestIdKey).(string)

		var requestBody string
		requestBody = `""`
		body, errJsonRequestBody := io.ReadAll(c.Request().Body)
		if errJsonRequestBody != nil {
			helpers.PrintLogToTerminal(errJsonRequestBody, requestId)
			errorMessages = helpers.ToErrorMessages("internal server error")
			err = errJsonRequestBody
			httpResponseCode = http.StatusInternalServerError
		}
		if len(body) == 0 {
			errLenBody := errors.New("json len body equal to 0")
			helpers.PrintLogToTerminal(errLenBody, requestId)
			err = errLenBody
			helpers.ToErrorMessages("internal server error")
		}
		jsonRequestBodyMap := make(map[string]interface{})
		errJsonRequestBodyMap := json.Unmarshal(body, &jsonRequestBodyMap)
		if errJsonRequestBodyMap != nil {
			helpers.PrintLogToTerminal(errJsonRequestBodyMap, requestId)
			err = errJsonRequestBodyMap
			errorMessages = helpers.ToErrorMessages("internal server error")
		}
		c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

		if jsonRequestBodyMap != nil {
			if c.Request().URL.Path == "/api/v1/users/register" {
				delete(jsonRequestBodyMap, "password")
				delete(jsonRequestBodyMap, "confirmpassword")
			} else if c.Request().URL.Path == "/api/v1/users/login" {
				delete(jsonRequestBodyMap, "password")
			}

			jsonRequestBodyByte, errJsonRequestBodyByte := json.Marshal(jsonRequestBodyMap)
			if errJsonRequestBodyByte != nil {
				helpers.PrintLogToTerminal(errJsonRequestBodyByte, requestId)
				errorMessages = helpers.ToErrorMessages("internal server error")
				httpResponseCode = http.StatusInternalServerError
				err = errJsonRequestBodyByte
			}
			requestBody = string(jsonRequestBodyByte)
		}

		host := c.Request().Host
		protocol := ""
		if c.Request().TLS == nil {
			protocol = "http"
		} else {
			protocol = "https"
		}
		urlPath := c.Request().URL.Path
		userAgent := c.Request().Header.Get("User-Agent")
		remoteAddr := c.Request().RemoteAddr
		forwardedFor := c.Request().Header.Get("X-Forwarded-For")

		requestLog := `{"requestTime": "` + datetimeNowRequest.String() + `", "app": "project-backend", "method": "` + requestMethod + `","requestId":"` + requestId + `","host": "` + host + `","urlPath":"` + urlPath + `","protocol":"` + protocol + `","body": ` + requestBody + `, "userAgent": "` + userAgent + `", "remoteAddr": "` + remoteAddr + `", "forwardedFor": "` + forwardedFor + `"}`
		fmt.Println(requestLog)
		if err != nil {
			return c.JSON(httpResponseCode, helpers.Response{Data: nil, Errors: errorMessages})
		}

		// so i can catch the response body
		resBody := new(bytes.Buffer)
		mw := io.MultiWriter(c.Response().Writer, resBody)
		writer := &responseBodyWriter{
			Writer:         mw,
			ResponseWriter: c.Response().Writer,
		}
		c.Response().Writer = writer

		err = next(c)
		if err != nil {
			c.Error(err)
		}

		responseBody := resBody.String()
		responseStatus := writer.status
		log := `{"responseTime": "` + time.Now().String() + `", "app": "project-backend", "requestId": "` + requestId + `", "responseStatus": ` + strconv.Itoa(responseStatus) + `, "response": ` + responseBody + `}`
		fmt.Println(log)
		return nil
	}
}

func PrintRequestResponseLogWithNoRequestBody(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		datetimeNowRequest := time.Now()
		requestMethod := c.Request().Method
		var err error
		var errorMessages []helpers.ErrorMessage
		var httpResponseCode int

		// why put this, because i wanna log request when there is no request id
		// requestId := c.Request().Header.Get("X-REQUEST-ID")
		// if requestId == "" {
		// 	errRequestId := errors.New("cannot find requestId")
		// 	helpers.PrintLogToTerminal(errRequestId, "requestId")
		// 	httpResponseCode = http.StatusBadRequest
		// 	errorMessages = helpers.ToErrorMessages("cannot find requestId")
		// 	err = errRequestId
		// }
		// ctx := context.WithValue(c.Request().Context(), RequestIdKey, requestId)
		// c.SetRequest(c.Request().WithContext(ctx))
		requestId := c.Request().Context().Value(RequestIdKey).(string)

		requestBody := `""`

		host := c.Request().Host
		protocol := ""
		if c.Request().TLS == nil {
			protocol = "http"
		} else {
			protocol = "https"
		}
		urlPath := c.Request().URL.Path
		userAgent := c.Request().Header.Get("User-Agent")
		remoteAddr := c.Request().RemoteAddr
		forwardedFor := c.Request().Header.Get("X-Forwarded-For")

		requestLog := `{"requestTime": "` + datetimeNowRequest.String() + `", "app": "project-backend", "method": "` + requestMethod + `","requestId":"` + requestId + `","host": "` + host + `","urlPath":"` + urlPath + `","protocol":"` + protocol + `","body": ` + requestBody + `, "userAgent": "` + userAgent + `", "remoteAddr": "` + remoteAddr + `", "forwardedFor": "` + forwardedFor + `"}`
		fmt.Println(requestLog)
		if err != nil {
			return c.JSON(httpResponseCode, helpers.Response{Data: nil, Errors: errorMessages})
		}

		// so i can catch the response body
		resBody := new(bytes.Buffer)
		mw := io.MultiWriter(c.Response().Writer, resBody)
		writer := &responseBodyWriter{Writer: mw, ResponseWriter: c.Response().Writer}
		c.Response().Writer = writer

		err = next(c)
		if err != nil {
			c.Error(err)
		}

		responseBody := resBody.String()
		log := `{"responseTime": "` + time.Now().String() + `", "app": "project-backend", "requestId": "` + requestId + `", "response": ` + responseBody + `}`
		fmt.Println(log)
		return nil
	}
}
