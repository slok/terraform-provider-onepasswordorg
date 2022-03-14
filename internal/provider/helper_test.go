package provider_test

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

// getFakeRepoTmpFile returns a temp file that can be used for the fake
// repository storage.
// It returns the file path and a delete file function
func getFakeRepoTmpFile(prefix string) (path string, delete func()) {
	// Prepare fake storage.
	f, err := os.CreateTemp("", prefix)
	if err != nil {
		panic(err)
	}
	return f.Name(), func() { _ = os.Remove(f.Name()) }
}

// assertUserOnFakeStorage is a helper to assert the expected user is stored on the fake
// repository.
func assertUserOnFakeStorage(t *testing.T, expUser *model.User) resource.TestCheckFunc {
	assert := assert.New(t)

	return resource.TestCheckFunc(func(s *terraform.State) error {
		// Get fake repo.
		repo := getFakeRepository(t)

		// Check user.
		gotUser, err := repo.GetUserByID(context.TODO(), expUser.ID)
		assert.NoError(err)
		assert.Equal(expUser, gotUser)
		return nil
	})
}

// assertUserDeletedOnFakeStorage is a helper to assert the expected user ID is not stored on
// the fake repository.
func assertUserDeletedOnFakeStorage(t *testing.T, userID string) resource.TestCheckFunc {
	assert := assert.New(t)

	return resource.TestCheckFunc(func(s *terraform.State) error {
		// Get fake repo.
		repo := getFakeRepository(t)

		// Check user is missing.
		_, err := repo.GetUserByID(context.TODO(), userID)
		assert.Error(err)
		return nil
	})
}

// assertGroupOnFakeStorage is a helper to assert the expected group is stored on the fake
// repository.
func assertGroupOnFakeStorage(t *testing.T, expGroup *model.Group) resource.TestCheckFunc {
	assert := assert.New(t)

	return resource.TestCheckFunc(func(s *terraform.State) error {
		// Get fake repo.
		repo := getFakeRepository(t)

		// Check user.
		gotGroup, err := repo.GetGroupByID(context.TODO(), expGroup.ID)
		assert.NoError(err)
		assert.Equal(expGroup, gotGroup)
		return nil
	})
}

// assertGroupDeletedOnFakeStorage is a helper to assert the expected user ID is not stored on
// the fake repository.
func assertGroupDeletedOnFakeStorage(t *testing.T, groupID string) resource.TestCheckFunc {
	assert := assert.New(t)

	return resource.TestCheckFunc(func(s *terraform.State) error {
		// Get fake repo.
		repo := getFakeRepository(t)

		// Check user is missing.
		_, err := repo.GetGroupByID(context.TODO(), groupID)
		assert.Error(err)
		return nil
	})
}
