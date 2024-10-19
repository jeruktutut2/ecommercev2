package services_test

import (
	"backend-golang/commons/helpers"
	"backend-golang/commons/middlewares"
	"backend-golang/commons/setups"
	"backend-golang/commons/utils"
	"backend-golang/features/users/login/models"
	"backend-golang/features/users/login/repositories"
	"backend-golang/features/users/login/services"
	"backend-golang/tests/initialize"
	"context"
	"net/http"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type LoginServiceTestSuite struct {
	suite.Suite
	ctx                      context.Context
	postgresUtil             utils.PostgresUtil
	redisUtil                utils.RedisUtil
	loginRequest             models.LoginRequest
	validate                 *validator.Validate
	userRepository           repositories.UserRepository
	userPermissionRepository repositories.UserPermissionRepository
	uuidHelper               helpers.UuidHelper
	residHelper              helpers.RedisHelper
	loginService             services.LoginService
}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginServiceTestSuite))
}

func (sut *LoginServiceTestSuite) SetupSuite() {
	sut.T().Log("SetupSuite")
	sut.postgresUtil = utils.NewPostgresConnection()
	sut.redisUtil = utils.NewRedisConnection()
	sut.validate = validator.New()
	setups.UsernameValidator(sut.validate)
	setups.PasswordValidator(sut.validate)
	setups.TelephoneValidator(sut.validate)
	sut.userRepository = repositories.NewUserRepository()
	sut.userPermissionRepository = repositories.NewUserPermissinoRepository()
	sut.uuidHelper = helpers.NewUuidHelper()
	sut.residHelper = helpers.NewRedisHelper()
	sut.loginService = services.NewLoginService(sut.postgresUtil, sut.redisUtil, sut.validate, sut.userRepository, sut.userPermissionRepository, sut.uuidHelper, sut.residHelper)
}

func (sut *LoginServiceTestSuite) SetupTest() {
	sut.T().Log("SetupTest")
	sut.ctx = context.WithValue(context.Background(), middlewares.RequestIdKey, uuid.New().String())
	sut.loginRequest = models.LoginRequest{
		Email:    "email@email.com",
		Password: "password@A1",
	}
}

func (sut *LoginServiceTestSuite) BeforeTest(suiteName, testName string) {
	sut.T().Log("BeforeTest: " + suiteName + " " + testName)
}

func (sut *LoginServiceTestSuite) Test1LoginValidationError() {
	sut.T().Log("Test1LoginValidationError")
	initialize.DropTableUserPermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTablePermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTableUser(sut.postgresUtil.GetPool(), sut.ctx)
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

func (sut *LoginServiceTestSuite) Test2LoginUserRepositoryFindByEmailInternalServerError() {
	sut.T().Log("Test2LoginUserRepositoryFindByEmailInternalServerError")
	initialize.DropTableUserPermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTablePermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	sessionId, httpCode, response := sut.loginService.Login(sut.ctx, sut.loginRequest)
	sut.Equal(sessionId, "")
	sut.Equal(httpCode, http.StatusInternalServerError)
	sut.Equal(response.Data, nil)
	errorMessages, _ := response.Errors.([]helpers.ErrorMessage)
	sut.Equal(errorMessages[0].Field, "message")
	sut.Equal(errorMessages[0].Message, "internal server error")
}

func (sut *LoginServiceTestSuite) Test3LoginUserRepositoryFindByEmailBadRequestWrongEmailOrPasswordError() {
	sut.T().Log("Test3LoginUserRepositoryFindByEmailBadRequestWrongEmailOrPasswordError")
	initialize.DropTableUserPermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTablePermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.CreateTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	sessionId, httpCode, response := sut.loginService.Login(sut.ctx, sut.loginRequest)
	sut.Equal(sessionId, "")
	sut.Equal(httpCode, http.StatusBadRequest)
	sut.Equal(response.Data, nil)
	errorMessages, _ := response.Errors.([]helpers.ErrorMessage)
	sut.Equal(errorMessages[0].Field, "message")
	sut.Equal(errorMessages[0].Message, "wrong email or password")
}

func (sut *LoginServiceTestSuite) Test4LoginBcryptCompareHashAndPasswordBadRequestWrongEmailOrPasswordError() {
	sut.T().Log("Test4LoginBcryptCompareHashAndPasswordBadRequestWrongEmailOrPasswordError")
	initialize.DropTableUserPermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTablePermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.CreateTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.CreateDataUser(sut.postgresUtil.GetPool(), sut.ctx)
	sut.loginRequest.Password = "password@A1-"
	sessionId, httpCode, response := sut.loginService.Login(sut.ctx, sut.loginRequest)
	sut.Equal(sessionId, "")
	sut.Equal(httpCode, http.StatusBadRequest)
	sut.Equal(response.Data, nil)
	errorMessages, _ := response.Errors.([]helpers.ErrorMessage)
	sut.Equal(errorMessages[0].Field, "message")
	sut.Equal(errorMessages[0].Message, "wrong email or password")
}

func (sut *LoginServiceTestSuite) Test5LoginUserPermissionRepositoryFindByUserIdInternalServerError() {
	sut.T().Log("Test5LoginUserPermissionRepositoryFindByUserIdInternalServerError")
	initialize.DropTableUserPermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTablePermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.CreateTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.CreateDataUser(sut.postgresUtil.GetPool(), sut.ctx)
	sessionId, httpCode, response := sut.loginService.Login(sut.ctx, sut.loginRequest)
	sut.Equal(sessionId, "")
	sut.Equal(httpCode, http.StatusInternalServerError)
	sut.Equal(response.Data, nil)
	errorMessages, _ := response.Errors.([]helpers.ErrorMessage)
	sut.Equal(errorMessages[0].Field, "message")
	sut.Equal(errorMessages[0].Message, "internal server error")
}

func (sut *LoginServiceTestSuite) Test6LoginSuccess() {
	sut.T().Log("Test6LoginSuccess")
	initialize.DropTableUserPermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTablePermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.CreateTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.CreateTablePermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.CreateTableUserPermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.CreateDataUser(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.CreateDataPermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.CreateDataUserPermission(sut.postgresUtil.GetPool(), sut.ctx)
	sessionId, httpCode, response := sut.loginService.Login(sut.ctx, sut.loginRequest)
	sut.NotEqual(sessionId, "")
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
