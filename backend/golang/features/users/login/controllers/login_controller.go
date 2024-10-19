package controllers

import (
	"backend-golang/commons/helpers"
	"backend-golang/features/users/login/models"
	"backend-golang/features/users/login/services"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type LoginController interface {
	Login(c echo.Context) error
}

type LoginControllerImplementation struct {
	LoginService services.LoginService
}

func NewLoginController(loginService services.LoginService) LoginController {
	return &LoginControllerImplementation{
		LoginService: loginService,
	}
}

func (controller *LoginControllerImplementation) Login(c echo.Context) error {
	var loginRequest models.LoginRequest
	err := c.Bind(&loginRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Response{Data: nil, Errors: helpers.ToErrorMessages(err.Error())})
	}
	sessionId, httpCode, response := controller.LoginService.Login(c.Request().Context(), loginRequest)

	secure, err := strconv.ParseBool(os.Getenv("ECOMMERCEV2_COOKIE_SECURE"))
	if err != nil {
		httpCode, response = helpers.ToResponseInternalServerError()
		return c.JSON(httpCode, response)
	}
	cookie := &http.Cookie{
		Name:     "sessionId",
		Value:    sessionId,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Domain:   os.Getenv("ECOMMERCEV2_COOKIE_DOMAIN"),
	}
	c.SetCookie(cookie)
	return c.JSON(httpCode, response)
}
