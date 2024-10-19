package services

import (
	"backend-golang/commons/helpers"
	"backend-golang/commons/utils"
	"backend-golang/features/users/login/models"
	"backend-golang/features/users/login/repositories"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"backend-golang/commons/middlewares"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type LoginService interface {
	Login(ctx context.Context, loginRequest models.LoginRequest) (sessionId string, httpCode int, response helpers.Response)
}

type LoginServiceImplementation struct {
	PostgresUtil             utils.PostgresUtil
	RedisUtil                utils.RedisUtil
	Validate                 *validator.Validate
	UserRepository           repositories.UserRepository
	UserPermissionRepository repositories.UserPermissionRepository
	UuidHelper               helpers.UuidHelper
	RedisHelper              helpers.RedisHelper
}

func NewLoginService(postgresUtil utils.PostgresUtil, redisUtil utils.RedisUtil, validate *validator.Validate, userRepository repositories.UserRepository, userPermissionRepository repositories.UserPermissionRepository, uuidHelper helpers.UuidHelper, redisHelper helpers.RedisHelper) LoginService {
	return &LoginServiceImplementation{
		PostgresUtil:             postgresUtil,
		RedisUtil:                redisUtil,
		Validate:                 validate,
		UserRepository:           userRepository,
		UserPermissionRepository: userPermissionRepository,
		UuidHelper:               uuidHelper,
		RedisHelper:              redisHelper,
	}
}

func (service *LoginServiceImplementation) Login(ctx context.Context, loginRequest models.LoginRequest) (sessionId string, httpCode int, response helpers.Response) {
	requestId := ctx.Value(middlewares.RequestIdKey).(string)
	var err error
	err = service.Validate.Struct(loginRequest)
	if err != nil {
		validationResult := helpers.GetValidatorError(err, loginRequest)
		if validationResult != nil {
			httpCode, response = helpers.ToResponseRequestValidation(requestId, validationResult)
			return
		}
	}

	user, err := service.UserRepository.FindByEmail(service.PostgresUtil.GetPool(), ctx, loginRequest.Email)
	if err != nil && err != pgx.ErrNoRows {
		httpCode, response = helpers.ToResponseCheckError(err, requestId)
		return
	} else if err != nil && err == pgx.ErrNoRows {
		err = errors.New("wrong email or password")
		httpCode, response = helpers.ToResponseError(err, requestId, http.StatusBadRequest, "wrong email or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(loginRequest.Password))
	if err != nil {
		err = errors.New("wrong email or password")
		httpCode, response = helpers.ToResponseError(err, requestId, http.StatusBadRequest, "wrong email or password")
		return
	}

	userPermissions, err := service.UserPermissionRepository.FindByUserId(service.PostgresUtil.GetPool(), ctx, user.Id.Int32)
	if err != nil {
		httpCode, response = helpers.ToResponseCheckError(err, requestId)
		return
	}
	var idPermissions []int32
	for _, userPermission := range userPermissions {
		idPermissions = append(idPermissions, userPermission.PermissionId.Int32)
	}

	sessionId = service.UuidHelper.String()
	sessionValue := make(map[string]interface{})
	sessionValue["id"] = user.Id.Int32
	sessionValue["username"] = user.Username.String
	sessionValue["email"] = user.Email.String
	sessionValue["idPermissions"] = idPermissions
	sessionByte, err := json.Marshal(sessionValue)
	if err != nil {
		httpCode, response = helpers.ToResponseCheckError(err, requestId)
		return
	}
	session := string(sessionByte)

	_, err = service.RedisHelper.Set(service.RedisUtil.GetClient(), ctx, sessionId, session, 0)
	if err != nil && err != redis.Nil {
		httpCode, response = helpers.ToResponseCheckError(err, requestId)
		return
	}

	httpCode = http.StatusOK
	responseMessage := helpers.ResponseMessage{
		Message: "successfully login",
	}
	response = helpers.Response{
		Data:   responseMessage,
		Errors: nil,
	}
	return
}
