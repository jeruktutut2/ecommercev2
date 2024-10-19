package setups

import (
	"backend-golang/commons/helpers"
	"backend-golang/commons/middlewares"
	"backend-golang/commons/utils"
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	loginroutes "backend-golang/features/users/login/routes"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func SetEcho(postgresUtil utils.PostgresUtil, redisUtil utils.RedisUtil, validate *validator.Validate, uuidHelper helpers.UuidHelper, redisHelper helpers.RedisHelper) (e *echo.Echo) {
	e = echo.New()
	e.Use(echomiddleware.Recover())
	e.Use(middlewares.SetRequestId)
	e.HTTPErrorHandler = CustomHTTPErrorHandler
	loginroutes.LoginRoute(e, postgresUtil, redisUtil, validate, uuidHelper, redisHelper)
	return
}

func StartEcho(e *echo.Echo) {
	host := os.Getenv("ECOMMERCEV2_ECHO_HOST")
	go func() {
		if err := e.Start(host); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()
	println(time.Now().String(), "echo: started at", host)
}

func StopEcho(e *echo.Echo) {
	host := os.Getenv("ECOMMERCEV2_ECHO_HOST")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	cancel()
	println(time.Now().String(), "echo: shutdown properly at", host)
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	requestId := c.Request().Context().Value(middlewares.RequestIdKey).(string)
	helpers.PrintLogToTerminal(err, requestId)
	he, ok := err.(*echo.HTTPError)
	if !ok {
		err = errors.New("cannot convert error to echo.HTTPError")
		helpers.PrintLogToTerminal(err, requestId)
		httpCode, response := helpers.ToResponseInternalServerError()
		c.JSON(httpCode, response)
		return
	}

	var message string
	if he.Code == http.StatusNotFound {
		message = "not found"
	} else if he.Code == http.StatusMethodNotAllowed {
		message = "method not allowed"
	} else {
		message = "internal server error"
	}
	errorMessages := helpers.ToErrorMessages(message)
	response := helpers.Response{
		Data:   nil,
		Errors: errorMessages,
	}
	c.JSON(he.Code, response)
}
