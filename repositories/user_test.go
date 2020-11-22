package repositories

import (
	"testing"

	"github.com/ksungcaya/todo-echo/models"
	"github.com/ksungcaya/todo-echo/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func newUserModel() *models.User {
	return &models.User{
		Name:     "John Doe",
		Username: "jdoe",
		Email:    "jdoe@example.com",
		Password: "secret",
	}
}

type UserRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo UserRepository
}

// Load test env and Refresh db
func (suite *UserRepositoryTestSuite) SetupTest() {
	db, _ := test.InitTestDB()
	suite.db = db
	suite.repo = NewUserRepository(db)
}

func (suite *UserRepositoryTestSuite) TestCreate() {
	assert := assert.New(suite.T())

	user := newUserModel()
	err := suite.repo.Create(user)

	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	var u models.User
	suite.db.Where("username = ?", user.Username).First(&u)

	assert.NotEmpty(u.CreatedAt)
	assert.Equal(user.Username, u.Username)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
