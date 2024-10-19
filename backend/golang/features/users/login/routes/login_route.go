package routes

import (
	"backend-golang/commons/helpers"
	"backend-golang/commons/middlewares"
	"backend-golang/commons/utils"
	"backend-golang/features/users/login/controllers"
	"backend-golang/features/users/login/repositories"
	"backend-golang/features/users/login/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func LoginRoute(e *echo.Echo, postgresUtil utils.PostgresUtil, redisUtil utils.RedisUtil, validate *validator.Validate, uuidHelper helpers.UuidHelper, redisHelper helpers.RedisHelper) {
	userRepository := repositories.NewUserRepository()
	userPermissionRepository := repositories.NewUserPermissinoRepository()
	loginService := services.NewLoginService(postgresUtil, redisUtil, validate, userRepository, userPermissionRepository, uuidHelper, redisHelper)
	loginController := controllers.NewLoginController(loginService)
	e.POST("/api/v1/users/login", loginController.Login, middlewares.PrintRequestResponseLog)
}
