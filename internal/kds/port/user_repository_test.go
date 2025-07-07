package port_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/datasource"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
	"github.com/leftovers-2025/kds_backend/internal/kds/repository/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository(t *testing.T) {
	repositories := createRepositories(t)
	for i := range repositories {
		repo := repositories[i]
		t.Run(repo.Name, func(t *testing.T) {
			t.Run("testCreateUser", func(t *testing.T) {
				testCreateUser(t, repo.Repository)
			})
			t.Run("testFindByGoogleId", func(t *testing.T) {
				testFindByGoogleId(t, repo.Repository)
			})
		})
	}
}

func testCreateUser(t *testing.T, repository port.UserRepository) {
	user := newValidUser(t)

	err := repository.Create(user)
	require.NoError(t, err)

	assertInRepository(t, repository, user)
}

func testFindByGoogleId(t *testing.T, repository port.UserRepository) {
	user := newValidUser(t)

	err := repository.Create(user)
	require.NoError(t, err)

	userInRepo, err := repository.FindByGoogleId(user.GoogleId())
	require.NoError(t, err)

	assert.Equal(t, user.Id(), userInRepo.Id())
}

func assertInRepository(t *testing.T, repository port.UserRepository, user *entity.User) {
	require.NotNil(t, user)

	userInRepo, err := repository.FindById(user.Id())
	require.NoError(t, err)

	assert.Equal(t, user.Id(), userInRepo.Id())
}

func newValidUser(t *testing.T) *entity.User {
	user, err := entity.NewUser(
		uuid.New(),
		uuid.NewString(),
		"validuser@example.com",
		uuid.NewString(),
		time.Now(),
		time.Now(),
	)
	require.NoError(t, err)

	return user
}

type UserRepositoryWithName struct {
	Name       string
	Repository port.UserRepository
}

func createRepositories(t *testing.T) []UserRepositoryWithName {
	return []UserRepositoryWithName{
		{
			Name:       "MySql",
			Repository: newMySqlRepository(t),
		},
	}
}

func newMySqlRepository(t *testing.T) port.UserRepository {
	repo := mysql.NewMySqlUserRepository(datasource.GetMySqlConnection())
	require.NotNil(t, repo)
	return repo
}
