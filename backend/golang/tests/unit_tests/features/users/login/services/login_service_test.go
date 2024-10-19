package services_test

import (
	"backend-golang/commons/helpers"
	"backend-golang/commons/middlewares"
	"backend-golang/commons/setups"
	"backend-golang/features/users/login/models"
	"backend-golang/features/users/login/services"
	mockhelpers "backend-golang/tests/unit_tests/commons/helpers/mocks"
	mockutils "backend-golang/tests/unit_tests/commons/utils/mocks"
	mockrepositories "backend-golang/tests/unit_tests/features/users/login/mocks/repositories"
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type LoginServiceTestSuite struct {
	suite.Suite
	ctx                          context.Context
	loginRequest                 models.LoginRequest
	postgresUtilMock             *mockutils.PostgresUtilMock
	redisUtilMock                *mockutils.RedisUtilMock
	validate                     *validator.Validate
	userRepositoryMock           *mockrepositories.UserRepositoryMock
	userPermissionRepositoryMock *mockrepositories.UserPermissionRepositoryMock
	uuidHelperMock               *mockhelpers.UuidHelperMock
	redisHelperMock              *mockhelpers.RedisHelperMock
	client                       *redis.Client
	pool                         *pgxpool.Pool
	errTimeout                   error
	errInternalServer            error
	user                         models.User
	sessionId                    string
	loginService                 services.LoginService
}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginServiceTestSuite))
}

func (sut *LoginServiceTestSuite) SetupSuite() {
	sut.T().Log("SetupSuite")
	sut.ctx = context.WithValue(context.Background(), middlewares.RequestIdKey, uuid.New().String())
	sut.client = &redis.Client{}
	sut.pool = &pgxpool.Pool{}
	sut.errTimeout = context.Canceled
	sut.errInternalServer = errors.New("internal server error")
	sut.sessionId = "sessionId"
}

func (sut *LoginServiceTestSuite) SetupTest() {
	sut.T().Log("SetupTest")
	sut.loginRequest = models.LoginRequest{
		Email:    "email@email.com",
		Password: "password@A1",
	}
	sut.user = models.User{
		Id:        pgtype.Int4{Valid: true, Int32: 1},
		Username:  pgtype.Text{Valid: true, String: "username"},
		Email:     pgtype.Text{Valid: true, String: "email@email.com"},
		Password:  pgtype.Text{Valid: true, String: "$2a$10$MvEM5qcQFk39jC/3fYzJzOIy7M/xQiGv/PAkkoarCMgsx/rO0UaPG"},
		CreatedAt: pgtype.Int8{Valid: true, Int64: 1719496855216},
	}
	sut.postgresUtilMock = new(mockutils.PostgresUtilMock)
	sut.redisUtilMock = new(mockutils.RedisUtilMock)
	sut.validate = validator.New()
	setups.UsernameValidator(sut.validate)
	setups.PasswordValidator(sut.validate)
	setups.TelephoneValidator(sut.validate)
	sut.userRepositoryMock = new(mockrepositories.UserRepositoryMock)
	sut.userPermissionRepositoryMock = new(mockrepositories.UserPermissionRepositoryMock)
	sut.uuidHelperMock = new(mockhelpers.UuidHelperMock)
	sut.redisHelperMock = new(mockhelpers.RedisHelperMock)
	sut.loginService = services.NewLoginService(sut.postgresUtilMock, sut.redisUtilMock, sut.validate, sut.userRepositoryMock, sut.userPermissionRepositoryMock, sut.uuidHelperMock, sut.redisHelperMock)
}

func (sut *LoginServiceTestSuite) BeforeTest(suiteName, testName string) {
	sut.T().Log("BeforeTest: " + suiteName + " " + testName)
}

func (sut *LoginServiceTestSuite) Test01LoginRedisRepositoryDelWithSessionIdTimeoutError() {
	sut.T().Log("Test01LoginRedisRepositoryDelWithSessionIdTimeoutError")
	sut.loginRequest = models.LoginRequest{}
	sessionId, httpCode, response := sut.loginService.Login(sut.ctx, sut.loginRequest)
	sut.Equal(sessionId, "")
	sut.Equal(httpCode, http.StatusBadRequest)
	sut.Equal(response.Data, nil)
	errorMessages, _ := response.Errors.([]helpers.ErrorMessage)
	sut.Equal(errorMessages[0].Field, "email")
	sut.Equal(errorMessages[0].Message, "is required")
	sut.Equal(errorMessages[1].Field, "password")
	sut.Equal(errorMessages[1].Message, "is required")
}

func (sut *LoginServiceTestSuite) Test02LoginUserRepositoryFindByEmailTimeoutError() {
	sut.T().Log("Test02LoginUserRepositoryFindByEmailTimeoutError")
	sut.postgresUtilMock.Mock.On("GetPool").Return(sut.pool)
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.pool, sut.ctx, sut.loginRequest.Email).Return(models.User{}, sut.errTimeout)
	sessionId, httpCode, response := sut.loginService.Login(sut.ctx, sut.loginRequest)
	sut.Equal(sessionId, "")
	sut.Equal(httpCode, http.StatusRequestTimeout)
	sut.Equal(response.Data, nil)
	errorMessages, _ := response.Errors.([]helpers.ErrorMessage)
	sut.Equal(errorMessages[0].Field, "message")
	sut.Equal(errorMessages[0].Message, "time out or user cancel the request")
}

func (sut *LoginServiceTestSuite) Test03LoginUserRepositoryFindByEmailInternalServerError() {
	sut.T().Log("Test03LoginUserRepositoryFindByEmailInternalServerError")
	sut.postgresUtilMock.Mock.On("GetPool").Return(sut.pool)
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.pool, sut.ctx, sut.loginRequest.Email).Return(models.User{}, sut.errInternalServer)
	sessionId, httpCode, response := sut.loginService.Login(sut.ctx, sut.loginRequest)
	sut.Equal(sessionId, "")
	sut.Equal(httpCode, http.StatusInternalServerError)
	sut.Equal(response.Data, nil)
	errorMessages, _ := response.Errors.([]helpers.ErrorMessage)
	sut.Equal(errorMessages[0].Field, "message")
	sut.Equal(errorMessages[0].Message, "internal server error")
}

func (sut *LoginServiceTestSuite) Test04LoginUserRepositoryFindByEmailBadRequestWrongEmailPassword() {
	sut.T().Log("Test04LoginUserRepositoryFindByEmailBadRequestWrongEmailPassword")
	sut.postgresUtilMock.Mock.On("GetPool").Return(sut.pool)
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.pool, sut.ctx, sut.loginRequest.Email).Return(models.User{}, pgx.ErrNoRows)
	sessionId, httpCode, response := sut.loginService.Login(sut.ctx, sut.loginRequest)
	sut.Equal(sessionId, "")
	sut.Equal(httpCode, http.StatusBadRequest)
	sut.Equal(response.Data, nil)
	errorMessages, _ := response.Errors.([]helpers.ErrorMessage)
	sut.Equal(errorMessages[0].Field, "message")
	sut.Equal(errorMessages[0].Message, "wrong email or password")
}

func (sut *LoginServiceTestSuite) Test05LoginBcryptCompareHashAndPasswordBadRequestWrongEmailPassword() {
	sut.T().Log("Test05LoginBcryptCompareHashAndPasswordBadRequestWrongEmailPassword")
	sut.postgresUtilMock.Mock.On("GetPool").Return(sut.pool)
	sut.user.Password.String = ""
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.pool, sut.ctx, sut.loginRequest.Email).Return(sut.user, nil)
	sessionId, httpCode, response := sut.loginService.Login(sut.ctx, sut.loginRequest)
	sut.Equal(sessionId, "")
	sut.Equal(httpCode, http.StatusBadRequest)
	sut.Equal(response.Data, nil)
	errorMessages, _ := response.Errors.([]helpers.ErrorMessage)
	sut.Equal(errorMessages[0].Field, "message")
	sut.Equal(errorMessages[0].Message, "wrong email or password")
}

func (sut *LoginServiceTestSuite) Test06LoginUserPermissionRepositoryFindByUserIdTimeoutError() {
	sut.T().Log("Test06LoginUserPermissionRepositoryFindByUserIdTimeoutError")
	sut.postgresUtilMock.Mock.On("GetPool").Return(sut.pool)
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.pool, sut.ctx, sut.loginRequest.Email).Return(sut.user, nil)
	sut.userPermissionRepositoryMock.Mock.On("FindByUserId", sut.pool, sut.ctx, sut.user.Id.Int32).Return([]models.UserPermission{}, sut.errTimeout)
	sessionId, httpCode, response := sut.loginService.Login(sut.ctx, sut.loginRequest)
	sut.Equal(sessionId, "")
	sut.Equal(httpCode, http.StatusRequestTimeout)
	sut.Equal(response.Data, nil)
	errorMessages, _ := response.Errors.([]helpers.ErrorMessage)
	sut.Equal(errorMessages[0].Field, "message")
	sut.Equal(errorMessages[0].Message, "time out or user cancel the request")
}

func (sut *LoginServiceTestSuite) Test07LoginUserPermissionRepositoryFindByUserIdInternalServerError() {
	sut.T().Log("Test07LoginUserPermissionRepositoryFindByUserIdInternalServerError")
	sut.postgresUtilMock.Mock.On("GetPool").Return(sut.pool)
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.pool, sut.ctx, sut.loginRequest.Email).Return(sut.user, nil)
	sut.userPermissionRepositoryMock.Mock.On("FindByUserId", sut.pool, sut.ctx, sut.user.Id.Int32).Return([]models.UserPermission{}, sut.errInternalServer)
	sessionId, httpCode, response := sut.loginService.Login(sut.ctx, sut.loginRequest)
	sut.Equal(sessionId, "")
	sut.Equal(httpCode, http.StatusInternalServerError)
	sut.Equal(response.Data, nil)
	errorMessages, _ := response.Errors.([]helpers.ErrorMessage)
	sut.Equal(errorMessages[0].Field, "message")
	sut.Equal(errorMessages[0].Message, "internal server error")
}

func (sut *LoginServiceTestSuite) Test08LoginRedisRepositorySetTimeoutError() {
	sut.T().Log("Test08LoginRedisRepositorySetTimeoutError")
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.postgresUtilMock.Mock.On("GetPool").Return(sut.pool)
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.pool, sut.ctx, sut.loginRequest.Email).Return(sut.user, nil)
	sut.userPermissionRepositoryMock.Mock.On("FindByUserId", sut.pool, sut.ctx, sut.user.Id.Int32).Return([]models.UserPermission{}, nil)
	sut.uuidHelperMock.Mock.On("String").Return(sut.sessionId)
	session := `{"email":"email@email.com","id":1,"idPermissions":null,"username":"username"}`
	sut.redisHelperMock.Mock.On("Set", sut.client, sut.ctx, sut.sessionId, session, time.Duration(0)).Return("", sut.errTimeout)
	sessionId, httpCode, response := sut.loginService.Login(sut.ctx, sut.loginRequest)
	sut.Equal(sessionId, sut.sessionId)
	sut.Equal(httpCode, http.StatusRequestTimeout)
	sut.Equal(response.Data, nil)
	errorMessages, _ := response.Errors.([]helpers.ErrorMessage)
	sut.Equal(errorMessages[0].Field, "message")
	sut.Equal(errorMessages[0].Message, "time out or user cancel the request")
}

func (sut *LoginServiceTestSuite) Test09LoginRedisRepositorySetInternalServerError() {
	sut.T().Log("Test09LoginRedisRepositorySetInternalServerError")
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.postgresUtilMock.Mock.On("GetPool").Return(sut.pool)
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.pool, sut.ctx, sut.loginRequest.Email).Return(sut.user, nil)
	sut.userPermissionRepositoryMock.Mock.On("FindByUserId", sut.pool, sut.ctx, sut.user.Id.Int32).Return([]models.UserPermission{}, nil)
	sut.uuidHelperMock.Mock.On("String").Return(sut.sessionId)
	session := `{"email":"email@email.com","id":1,"idPermissions":null,"username":"username"}`
	sut.redisHelperMock.Mock.On("Set", sut.client, sut.ctx, sut.sessionId, session, time.Duration(0)).Return("", sut.errInternalServer)
	sessionId, httpCode, response := sut.loginService.Login(sut.ctx, sut.loginRequest)
	sut.Equal(sessionId, sut.sessionId)
	sut.Equal(httpCode, http.StatusInternalServerError)
	sut.Equal(response.Data, nil)
	errorMessages, _ := response.Errors.([]helpers.ErrorMessage)
	sut.Equal(errorMessages[0].Field, "message")
	sut.Equal(errorMessages[0].Message, "internal server error")
}

func (sut *LoginServiceTestSuite) Test10LoginSuccess() {
	sut.T().Log("Test13LoginSuccess")
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.postgresUtilMock.Mock.On("GetPool").Return(sut.pool)
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.pool, sut.ctx, sut.loginRequest.Email).Return(sut.user, nil)
	sut.userPermissionRepositoryMock.Mock.On("FindByUserId", sut.pool, sut.ctx, sut.user.Id.Int32).Return([]models.UserPermission{}, nil)
	sut.uuidHelperMock.Mock.On("String").Return(sut.sessionId)
	session := `{"email":"email@email.com","id":1,"idPermissions":null,"username":"username"}`
	sut.redisHelperMock.Mock.On("Set", sut.client, sut.ctx, sut.sessionId, session, time.Duration(0)).Return("", nil)
	sessionId, httpCode, response := sut.loginService.Login(sut.ctx, sut.loginRequest)
	sut.Equal(sessionId, sut.sessionId)
	sut.Equal(httpCode, http.StatusOK)
	sut.Equal(response.Errors, nil)
	responseMessage, _ := response.Data.(helpers.ResponseMessage)
	sut.Equal(responseMessage.Message, "successfully login")
}

func (sut *LoginServiceTestSuite) AfterTest(suiteName, testName string) {
	sut.T().Log("AfterTest: " + suiteName + " " + testName)
}

func (sut *LoginServiceTestSuite) TearDownTest() {
	sut.T().Log("TearDownTest")
}

func (sut *LoginServiceTestSuite) TearDownSuite() {
	sut.T().Log("TearDownSuite")
}
