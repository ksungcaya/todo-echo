package repositories

import (
	"testing"

	"github.com/ksungcaya/todo-echo/models"
	"github.com/ksungcaya/todo-echo/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func existingUser() *models.User {
	return &models.User{
		Username: "janedoe",
		Name:     "Jane Doe",
		Email:    "janedoe@example.com",
		Password: "secret",
	}
}

type UserRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo UserRepository
	user *models.User
}

// Load test env and Refresh db
func (suite *UserRepositoryTestSuite) SetupTest() {
	db, _ := test.InitTestDB()
	// Truncate users table
	db.Where("1 = 1").Delete(&models.User{})

	suite.db = db
	suite.repo = NewUserRepository(db)

	user := existingUser()
	suite.db.Create(user)
	suite.user = user
}

func (suite *UserRepositoryTestSuite) TestByID() {
	assert := assert.New(suite.T())

	user := suite.repo.ByID(suite.user.ID)

	assert.NotEmpty(user)
	assert.Equal(suite.user.Email, user.Email)
}

func (suite *UserRepositoryTestSuite) TestByUsername() {
	assert := assert.New(suite.T())

	user := suite.repo.ByUsername(suite.user.Username)

	assert.NotEmpty(user)
	assert.Equal(suite.user.Email, user.Email)
}

func (suite *UserRepositoryTestSuite) TestByEmail() {
	assert := assert.New(suite.T())

	user := suite.repo.ByEmail(suite.user.Email)

	assert.NotEmpty(user)
	assert.Equal(suite.user.Username, user.Username)
}

func (suite *UserRepositoryTestSuite) TestCreate() {
	assert := assert.New(suite.T())

	user := &models.User{
		Username: "jdoe",
		Name:     "John Doe",
		Email:    "jdoe@example.com",
		Password: "secret",
	}
	err := suite.repo.Create(user)

	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	var u models.User
	suite.db.Where("username = ?", user.Username).First(&u)

	assert.NotEmpty(u.CreatedAt)
	assert.Equal(user.Username, u.Username)

	// assert password hashed
	assert.False(u.CheckPassword("wrongpassword"))
	assert.True(u.CheckPassword("secret"))
	assert.NotEqual("secret", u.Password)
}

func (suite *UserRepositoryTestSuite) TestUpdate() {
	assert := assert.New(suite.T())

	var u models.User
	suite.db.Where("username = ?", suite.user.Username).First(&u)

	// update
	u.Username = ""
	u.Name = "Mary Jane Doe"
	suite.repo.Update(&u)

	var updated models.User
	suite.db.Where("username = ?", suite.user.Username).First(&updated)

	assert.NotEmpty(updated.Username) // username should not be updated
	assert.NotEqual(suite.user.Name, updated.Name)
	assert.Equal("Mary Jane Doe", updated.Name)
	assert.True(updated.CheckPassword("secret")) // password should not be updated
	assert.Greater(updated.UpdatedAt.UnixNano(), u.CreatedAt.UnixNano())
}

func (suite *UserRepositoryTestSuite) TestDelete() {
	assert := assert.New(suite.T())

	var u models.User
	suite.db.Where("username = ?", suite.user.Username).First(&u)
	assert.NotEmpty(u)

	// delete
	suite.repo.Delete(u.ID)

	var deleted models.User
	suite.db.Where("username = ?", suite.user.Username).First(&deleted)

	assert.Empty(deleted)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
