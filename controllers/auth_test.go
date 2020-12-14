package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mocks "github.com/ksungcaya/todo-echo/mocks/repositories"
	"github.com/ksungcaya/todo-echo/models"
	"github.com/ksungcaya/todo-echo/requests"
	"github.com/ksungcaya/todo-echo/test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthControllerTestSuite struct {
	suite.Suite
	repo   *mocks.UserRepository
	auth   *AuthController
	server *echo.Echo
}

var (
	registerRequest = requests.RegisterRequest{
		Username: "alice",
		Name:     "Alice Wonder",
		Email:    "alice@realworld.io",
		Password: "secret",
	}

	loginRequest = requests.LoginRequest{
		Username: "alice",
		Password: "secret",
	}

	existingUser = &models.User{
		Username: registerRequest.Username,
		Name:     registerRequest.Name,
		Email:    registerRequest.Email,
		Password: "$2a$10$vFN7/BdlTDFcp1ndGQELtu4eRY6MtEccXJ3tUwfP4qAzMMfDaypBe", // secret
	}
)

// Setup auth
func (suite *AuthControllerTestSuite) SetupTest() {
	suite.repo = &mocks.UserRepository{}
	suite.auth = NewAuth(suite.repo)
	suite.server = echo.New()
}

func (suite *AuthControllerTestSuite) TestLoignInvalidPayload() {
	assert := assert.New(suite.T())

	request := httptest.NewRequest(echo.POST, "/auth/login", bytes.NewBufferString("invalid json format"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	context := suite.server.NewContext(request, response)

	assert.NoError(suite.auth.Login(context))

	if assert.Equal(http.StatusBadRequest, response.Code) {
		err := test.GetResponseErrors(response)
		assert.Equal(err["message"], "Invalid request payload")
	}
}

func (suite *AuthControllerTestSuite) TestLoginValidation() {
	assert := assert.New(suite.T())
	invalidPayload := `{
		"username": "",
		"password": ""
	}`

	request := httptest.NewRequest(echo.POST, "/auth/login", strings.NewReader(invalidPayload))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	context := suite.server.NewContext(request, response)

	assert.NoError(suite.auth.Register(context))

	if assert.Equal(http.StatusUnprocessableEntity, response.Code) {
		err := test.GetResponseErrors(response)
		assert.NotEmpty(err["username"])
		assert.NotEmpty(err["password"])
	}
}

func (suite *AuthControllerTestSuite) TestLoginWithUnregisteredUser() {
	assert := assert.New(suite.T())

	loginPayload, _ := json.Marshal(loginRequest)
	request := httptest.NewRequest(echo.POST, "/auth/login", strings.NewReader(string(loginPayload)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	context := suite.server.NewContext(request, response)

	suite.repo.On("ByUsername", loginRequest.Username).Return(nil)

	assert.NoError(suite.auth.Login(context))

	suite.repo.AssertCalled(suite.T(), "ByUsername", loginRequest.Username)
	if assert.Equal(http.StatusForbidden, response.Code) {
		err := test.GetResponseErrors(response)
		assert.Equal(err["message"], "Invalid username or password")
	}
}

func (suite *AuthControllerTestSuite) TestLoginWithInvalidPassword() {
	assert := assert.New(suite.T())

	invalidPasswordRequest := loginRequest
	invalidPasswordRequest.Password = "invalid-pass"
	loginPayload, _ := json.Marshal(invalidPasswordRequest)
	request := httptest.NewRequest(echo.POST, "/auth/login", strings.NewReader(string(loginPayload)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	context := suite.server.NewContext(request, response)

	suite.repo.On("ByUsername", loginRequest.Username).Return(existingUser)

	assert.NoError(suite.auth.Login(context))

	suite.repo.AssertCalled(suite.T(), "ByUsername", loginRequest.Username)
	if assert.Equal(http.StatusForbidden, response.Code) {
		err := test.GetResponseErrors(response)
		assert.Equal(err["message"], "Invalid username or password")
	}
}

func (suite *AuthControllerTestSuite) TestLoginSuccess() {
	assert := assert.New(suite.T())

	loginPayload, _ := json.Marshal(loginRequest)
	request := httptest.NewRequest(echo.POST, "/auth/login", strings.NewReader(string(loginPayload)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	context := suite.server.NewContext(request, response)

	suite.repo.On("ByUsername", loginRequest.Username).Return(existingUser)

	assert.NoError(suite.auth.Login(context))

	suite.repo.AssertCalled(suite.T(), "ByUsername", loginRequest.Username)
	if assert.Equal(http.StatusOK, response.Code) {
		data := test.GetResponseData(response)
		suite.T().Logf("\n%v", data)
		assert.NotEmpty(data["token"])
	}
}

func (suite *AuthControllerTestSuite) TestRegistrationInvalidPayload() {
	assert := assert.New(suite.T())

	request := httptest.NewRequest(echo.POST, "/auth/register", bytes.NewBufferString("invalid json format"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	context := suite.server.NewContext(request, response)

	assert.NoError(suite.auth.Register(context))

	if assert.Equal(http.StatusBadRequest, response.Code) {
		err := test.GetResponseErrors(response)
		assert.Equal(err["message"], "Invalid request payload")
	}
}

func (suite *AuthControllerTestSuite) TestRegistrationValidation() {
	assert := assert.New(suite.T())

	invalidPayload := `{
		"username": "",
		"name": "",
		"email": "alice",
		"password": ""
	}`

	request := httptest.NewRequest(echo.POST, "/auth/register", strings.NewReader(string(invalidPayload)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	context := suite.server.NewContext(request, response)

	assert.NoError(suite.auth.Register(context))

	if assert.Equal(http.StatusUnprocessableEntity, response.Code) {
		err := test.GetResponseErrors(response)
		assert.NotEmpty(err["username"])
		assert.NotEmpty(err["email"])
		assert.NotEmpty(err["name"])
		assert.NotEmpty(err["password"])

		// has email but invalid
		assert.Len(err["email"], 1)
	}
}

func (suite *AuthControllerTestSuite) TestRegistrationWithExistingEmail() {
	assert := assert.New(suite.T())

	registerPayload, _ := json.Marshal(registerRequest)
	request := httptest.NewRequest(echo.POST, "/auth/register", strings.NewReader(string(registerPayload)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	context := suite.server.NewContext(request, response)

	suite.repo.On("ByEmail", registerRequest.Email).Return(existingUser)

	assert.NoError(suite.auth.Register(context))

	suite.repo.AssertCalled(suite.T(), "ByEmail", registerRequest.Email)
	if assert.Equal(http.StatusUnprocessableEntity, response.Code) {
		err := test.GetResponseErrors(response)
		assert.NotEmpty(err["email"])
		assert.Contains(err["email"].([]interface{})[0], "exist")
	}
}

func (suite *AuthControllerTestSuite) TestRegistrationSuccess() {
	assert := assert.New(suite.T())

	registerPayload, _ := json.Marshal(registerRequest)
	request := httptest.NewRequest(echo.POST, "/auth/register", strings.NewReader(string(registerPayload)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	context := suite.server.NewContext(request, response)

	user := existingUser
	user.Password = "secret"

	suite.repo.On("ByEmail", registerRequest.Email).Return(nil)
	suite.repo.On("Create", user).Return(nil)

	assert.NoError(suite.auth.Register(context))

	suite.repo.AssertCalled(suite.T(), "Create", user)
	if assert.Equal(http.StatusOK, response.Code) {
		data := test.GetResponseData(response)
		assert.NotEmpty(data["name"])
		assert.NotEmpty(data["email"])
		assert.NotEmpty(data["username"])
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAuthControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}
