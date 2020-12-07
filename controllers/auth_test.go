package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mocks "github.com/ksungcaya/todo-echo/mocks/repositories"
	"github.com/ksungcaya/todo-echo/models"
	"github.com/ksungcaya/todo-echo/requests"
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
)

// Setup auth
func (suite *AuthControllerTestSuite) SetupTest() {
	suite.repo = &mocks.UserRepository{}
	suite.auth = NewAuth(suite.repo)
	suite.server = echo.New()
}

func (suite *AuthControllerTestSuite) TestRegistrationValidation() {
	assert := assert.New(suite.T())

	invalidPayload := `{
		"username": "",
		"name": "",
		"email": "alice",
		"password": ""
	}`

	request := httptest.NewRequest(echo.POST, "/auth/register", strings.NewReader(invalidPayload))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	context := suite.server.NewContext(request, response)

	assert.NoError(suite.auth.Register(context))

	if assert.Equal(http.StatusUnprocessableEntity, response.Code) {
		var data map[string]interface{}
		json.Unmarshal(response.Body.Bytes(), &data)
		err := data["errors"].(map[string]interface{})

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

	suite.repo.On("ByEmail", registerRequest.Email).Return(&models.User{
		Username: registerRequest.Username,
		Name:     registerRequest.Name,
		Email:    registerRequest.Email,
		Password: registerRequest.Password,
	})

	assert.NoError(suite.auth.Register(context))

	suite.repo.AssertCalled(suite.T(), "ByEmail", registerRequest.Email)
	if assert.Equal(http.StatusUnprocessableEntity, response.Code) {
		var data map[string]interface{}
		json.Unmarshal(response.Body.Bytes(), &data)
		err := data["errors"].(map[string]interface{})

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

	userPayload := &models.User{
		Username: registerRequest.Username,
		Name:     registerRequest.Name,
		Email:    registerRequest.Email,
		Password: registerRequest.Password,
	}

	suite.repo.On("ByEmail", registerRequest.Email).Return(nil)
	suite.repo.On("Create", userPayload).Return(nil)

	assert.NoError(suite.auth.Register(context))

	suite.repo.AssertCalled(suite.T(), "Create", userPayload)
	if assert.Equal(http.StatusOK, response.Code) {
		var res map[string]interface{}
		json.Unmarshal(response.Body.Bytes(), &res)
		data := res["data"].(map[string]interface{})
		suite.T().Logf("\n%v\n", data)

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
