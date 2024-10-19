package login_test

import (
	"backend-golang/commons/helpers"
	"backend-golang/commons/middlewares"
	"backend-golang/commons/setups"
	"backend-golang/commons/utils"
	"backend-golang/features/users/login/routes"
	"backend-golang/tests/initialize"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/suite"
)

type LoginTestSuite struct {
	suite.Suite
	ctx          context.Context
	postgresUtil utils.PostgresUtil
	redisUtil    utils.RedisUtil
	validate     *validator.Validate
	requestBody  string
	e            *echo.Echo
	uuidHelper   helpers.UuidHelper
	redisHelper  helpers.RedisHelper
}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}

func (sut *LoginTestSuite) SetupSuite() {
	sut.T().Log("SetupSuite")
	sut.postgresUtil = utils.NewPostgresConnection()
	sut.redisUtil = utils.NewRedisConnection()
	sut.validate = setups.SetValidator()
	sut.uuidHelper = helpers.NewUuidHelper()
	sut.redisHelper = helpers.NewRedisHelper()
	sut.e = echo.New()
	sut.e.Use(echomiddleware.Recover())
	sut.e.Use(middlewares.SetRequestId)
	sut.e.HTTPErrorHandler = setups.CustomHTTPErrorHandler
	routes.LoginRoute(sut.e, sut.postgresUtil, sut.redisUtil, sut.validate, sut.uuidHelper, sut.redisHelper)
}

func (sut *LoginTestSuite) SetupTest() {
	sut.T().Log("SetupTest")
	sut.ctx = context.Background()
	sut.requestBody = `{
		"email": "email@email.com",
		"password": "password@A1"
	}`
}

func (sut *LoginTestSuite) BeforeTest(suiteName, testName string) {
	sut.T().Log("BeforeTest: " + suiteName + " " + testName)
}

func (sut *LoginTestSuite) Test1LoginValidationError() {
	sut.T().Log("Test1LoginValidationError")
	initialize.DropTableUserPermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTablePermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	sut.requestBody = `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", strings.NewReader(sut.requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	sut.e.ServeHTTP(rec, req)
	response := rec.Result()
	sut.Equal(response.StatusCode, http.StatusBadRequest)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	fmt.Println()
	fmt.Println("responseBody:", responseBody)
	fmt.Println()
	sut.Equal(responseBody["data"], nil)
	errorsResponseBody := responseBody["errors"].([]interface{})
	errorMessage0, _ := errorsResponseBody[0].((map[string]interface{}))
	sut.Equal(errorMessage0["field"], "email")
	sut.Equal(errorMessage0["message"], "is required")
	errorMessage1, _ := errorsResponseBody[1].((map[string]interface{}))
	sut.Equal(errorMessage1["field"], "password")
	sut.Equal(errorMessage1["message"], "is required")
}

func (sut *LoginTestSuite) Test2LoginUserRepositoryFindByEmailInternalServerError() {
	sut.T().Log("Test2LoginUserRepositoryFindByEmailInternalServerError")
	initialize.DropTableUserPermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTablePermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", strings.NewReader(sut.requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	sut.e.ServeHTTP(rec, req)
	response := rec.Result()
	sut.Equal(response.StatusCode, http.StatusInternalServerError)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	sut.Equal(responseBody["data"], nil)
	errorsResponseBody := responseBody["errors"].([]interface{})
	errorMessage0, _ := errorsResponseBody[0].((map[string]interface{}))
	sut.Equal(errorMessage0["field"], "message")
	sut.Equal(errorMessage0["message"], "internal server error")
}

func (sut *LoginTestSuite) Test3LoginUserRepositoryFindByEmailBadRequestWrongEmailOrPasswordError() {
	sut.T().Log("Test3LoginUserRepositoryFindByEmailBadRequestWrongEmailOrPasswordError")
	initialize.DropTableUserPermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTablePermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.CreateTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", strings.NewReader(sut.requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	sut.e.ServeHTTP(rec, req)
	response := rec.Result()
	sut.Equal(response.StatusCode, http.StatusBadRequest)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	sut.Equal(responseBody["data"], nil)
	errorsResponseBody := responseBody["errors"].([]interface{})
	errorMessage0, _ := errorsResponseBody[0].((map[string]interface{}))
	sut.Equal(errorMessage0["field"], "message")
	sut.Equal(errorMessage0["message"], "wrong email or password")
}

func (sut *LoginTestSuite) Test4LoginBcryptCompareHashAndPasswordBadRequestWrongEmailOrPasswordError() {
	sut.T().Log("Test4LoginBcryptCompareHashAndPasswordBadRequestWrongEmailOrPasswordError")
	initialize.DropTableUserPermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTablePermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.CreateTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.CreateDataUser(sut.postgresUtil.GetPool(), sut.ctx)
	sut.requestBody = `{
 		"email": "email@email.com",
 		"password": "password@A1-"
 	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", strings.NewReader(sut.requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	sut.e.ServeHTTP(rec, req)
	response := rec.Result()
	sut.Equal(response.StatusCode, http.StatusBadRequest)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	sut.Equal(responseBody["data"], nil)
	errorsResponseBody := responseBody["errors"].([]interface{})
	errorMessage0, _ := errorsResponseBody[0].((map[string]interface{}))
	sut.Equal(errorMessage0["field"], "message")
	sut.Equal(errorMessage0["message"], "wrong email or password")
}

func (sut *LoginTestSuite) Test5LoginUserPermissionRepositoryFindByUserIdInternalServerError() {
	sut.T().Log("Test5LoginUserPermissionRepositoryFindByUserIdInternalServerError")
	initialize.DropTableUserPermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTablePermission(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.DropTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.CreateTableUser(sut.postgresUtil.GetPool(), sut.ctx)
	initialize.CreateDataUser(sut.postgresUtil.GetPool(), sut.ctx)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", strings.NewReader(sut.requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	sut.e.ServeHTTP(rec, req)
	response := rec.Result()
	sut.Equal(response.StatusCode, http.StatusInternalServerError)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	sut.Equal(responseBody["data"], nil)
	errorsResponseBody := responseBody["errors"].([]interface{})
	errorMessage0, _ := errorsResponseBody[0].((map[string]interface{}))
	sut.Equal(errorMessage0["field"], "message")
	sut.Equal(errorMessage0["message"], "internal server error")
}

func (sut *LoginTestSuite) Test6LoginSuccess() {
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
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", strings.NewReader(sut.requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	sut.e.ServeHTTP(rec, req)
	response := rec.Result()
	sut.Equal(response.StatusCode, http.StatusOK)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	sut.NotEqual(responseBody["data"], "")
	sut.Equal(responseBody["errors"], nil)
}

func (sut *LoginTestSuite) AfterTest(suiteName, testName string) {
	sut.T().Log("AfterTest: " + suiteName + " " + testName)
}

func (sut *LoginTestSuite) TearDownTest() {
	sut.T().Log("TearDownTest")
}

func (sut *LoginTestSuite) TearDownSuite() {
	sut.T().Log("TearDownSuite")
}
