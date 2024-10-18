package provider_test

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

// getFakeRepoTmpFile returns a temp file that can be used for the fake repository storage.
// It returns the file path and a delete file function.
func getFakeRepoTmpFile(prefix string) (path string, delete func()) {
	// Prepare fake storage.
	f, err := os.CreateTemp("", prefix)
	if err != nil {
		panic(err)
	}
	return f.Name(), func() { _ = os.Remove(f.Name()) }
}

func assertUserOnFakeStorage(t *testing.T, expUser *model.User) resource.TestCheckFunc {
	assert := assert.New(t)

	return resource.TestCheckFunc(func(s *terraform.State) error {
		repo := getFakeRepository(t)

		gotUser, err := repo.GetUserByID(context.TODO(), expUser.ID)
		assert.NoError(err)
		assert.Equal(expUser, gotUser)
		return nil
	})
}

func assertUserDeletedOnFakeStorage(t *testing.T, userID string) resource.TestCheckFunc {
	assert := assert.New(t)

	return resource.TestCheckFunc(func(s *terraform.State) error {
		repo := getFakeRepository(t)

		_, err := repo.GetUserByID(context.TODO(), userID)
		assert.Error(err)
		return nil
	})
}

func assertGroupOnFakeStorage(t *testing.T, expGroup *model.Group) resource.TestCheckFunc {
	assert := assert.New(t)

	return resource.TestCheckFunc(func(s *terraform.State) error {
		repo := getFakeRepository(t)

		gotGroup, err := repo.GetGroupByID(context.TODO(), expGroup.ID)
		assert.NoError(err)
		assert.Equal(expGroup, gotGroup)
		return nil
	})
}

func assertGroupDeletedOnFakeStorage(t *testing.T, groupID string) resource.TestCheckFunc {
	assert := assert.New(t)

	return resource.TestCheckFunc(func(s *terraform.State) error {
		repo := getFakeRepository(t)

		_, err := repo.GetGroupByID(context.TODO(), groupID)
		assert.Error(err)
		return nil
	})
}

func assertGroupMemberOnFakeStorage(t *testing.T, expMembership *model.Membership) resource.TestCheckFunc {
	assert := assert.New(t)

	return resource.TestCheckFunc(func(s *terraform.State) error {
		repo := getFakeRepository(t)

		gotMembership, err := repo.GetMembershipByID(context.TODO(), expMembership.GroupID, expMembership.UserID)
		assert.NoError(err)
		assert.Equal(expMembership, gotMembership)
		return nil
	})
}

func assertGroupMemberDeletedOnFakeStorage(t *testing.T, groupID, userID string) resource.TestCheckFunc {
	assert := assert.New(t)

	return resource.TestCheckFunc(func(s *terraform.State) error {
		repo := getFakeRepository(t)

		_, err := repo.GetMembershipByID(context.TODO(), groupID, userID)
		assert.Error(err)
		return nil
	})
}

func assertVaultOnFakeStorage(t *testing.T, expVault *model.Vault) resource.TestCheckFunc {
	assert := assert.New(t)

	return resource.TestCheckFunc(func(s *terraform.State) error {
		repo := getFakeRepository(t)

		gotVault, err := repo.GetVaultByID(context.TODO(), expVault.ID)
		assert.NoError(err)
		assert.Equal(expVault, gotVault)
		return nil
	})
}

func assertVaultDeletedOnFakeStorage(t *testing.T, vaultID string) resource.TestCheckFunc {
	assert := assert.New(t)

	return resource.TestCheckFunc(func(s *terraform.State) error {
		repo := getFakeRepository(t)

		_, err := repo.GetVaultByID(context.TODO(), vaultID)
		assert.Error(err)
		return nil
	})
}

func assertVaultGroupAccessOnFakeStorage(t *testing.T, exp *model.VaultGroupAccess) resource.TestCheckFunc {
	assert := assert.New(t)

	return resource.TestCheckFunc(func(s *terraform.State) error {
		repo := getFakeRepository(t)

		got, err := repo.GetVaultGroupAccessByID(context.TODO(), exp.VaultID, exp.GroupID)
		assert.NoError(err)
		assert.Equal(exp, got)
		return nil
	})
}

func assertVaultGroupAccessDeletedOnFakeStorage(t *testing.T, vaultID, groupID string) resource.TestCheckFunc {
	assert := assert.New(t)

	return resource.TestCheckFunc(func(s *terraform.State) error {
		repo := getFakeRepository(t)

		_, err := repo.GetVaultGroupAccessByID(context.TODO(), vaultID, groupID)
		assert.Error(err)
		return nil
	})
}

func assertVaultUserAccessOnFakeStorage(t *testing.T, exp *model.VaultUserAccess) resource.TestCheckFunc {
	assert := assert.New(t)

	return resource.TestCheckFunc(func(s *terraform.State) error {
		repo := getFakeRepository(t)

		got, err := repo.GetVaultUserAccessByID(context.TODO(), exp.VaultID, exp.UserID)
		assert.NoError(err)
		assert.Equal(exp, got)
		return nil
	})
}

func assertVaultUserAccessDeletedOnFakeStorage(t *testing.T, vaultID, userID string) resource.TestCheckFunc {
	assert := assert.New(t)

	return resource.TestCheckFunc(func(s *terraform.State) error {
		repo := getFakeRepository(t)

		_, err := repo.GetVaultUserAccessByID(context.TODO(), vaultID, userID)
		assert.Error(err)
		return nil
	})
}
