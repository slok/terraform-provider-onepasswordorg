package onepasswordcli_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/storage/onepasswordcli"
	"github.com/slok/terraform-provider-onepasswordorg/internal/storage/onepasswordcli/onepasswordclimock"
)

func TestRepositoryEnsureMembership(t *testing.T) {
	tests := map[string]struct {
		membership model.Membership
		mock       func(m *onepasswordclimock.OpCli)
		expErr     bool
	}{
		"Creating a membership correctly, should return the data with the ID.": {
			membership: model.Membership{UserID: "test-00", GroupID: "group-00", Role: model.MembershipRoleMember},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `add user test-00 group-00 --role member`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)
			},
		},

		"If the user has wants a role other than member it shoul be called twice.": {
			membership: model.Membership{UserID: "test-00", GroupID: "group-00", Role: model.MembershipRoleManager},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `add user test-00 group-00 --role manager`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Twice().Return("", "", nil)
			},
		},

		"Having an error while calling the create op CLI action, should fail.": {
			membership: model.Membership{UserID: "test-00", GroupID: "group-00", Role: model.MembershipRoleMember},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `add user test-00 group-00 --role member`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("something"))
			},
			expErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			mc := &onepasswordclimock.OpCli{}
			test.mock(mc)

			repo, err := onepasswordcli.NewRepository(mc)
			require.NoError(err)

			err = repo.EnsureMembership(context.TODO(), test.membership)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryGetMembershipByID(t *testing.T) {
	tests := map[string]struct {
		userID        string
		groupID       string
		mock          func(m *onepasswordclimock.OpCli)
		expMembership *model.Membership
		expErr        bool
	}{
		"Getting a member correctly, should return the group data.": {
			userID:  "test-user-00",
			groupID: "group-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `list users --group group-id`
				stdout := `[{"uuid":"test-user-00","firstName":"Test00","lastName":"","name":"Test00","email":"test0@slok.dev","avatar":"","state":"A","type":"R","role":"MANAGER"},{"uuid":"test-user-01","firstName":"Test01","lastName":"","name":"Tst01","email":"test01@slok.dev","avatar":"","state":"3","type":"R","role":"MEMBER"}]`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expMembership: &model.Membership{
				UserID:  "test-user-00",
				GroupID: "group-id",
				Role:    model.MembershipRoleManager,
			},
		},

		"Getting a missing member should fail.": {
			userID:  "test-user-02",
			groupID: "group-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `list users --group group-id`
				stdout := `[{"uuid":"test-user-00","firstName":"Test00","lastName":"","name":"Test00","email":"test0@slok.dev","avatar":"","state":"A","type":"R","role":"MANAGER"},{"uuid":"test-user-01","firstName":"Test01","lastName":"","name":"Tst01","email":"test01@slok.dev","avatar":"","state":"3","type":"R","role":"MEMBER"}]`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expErr: true,
		},

		"Having an error while calling the op CLI, should fail.": {
			userID:  "test-id",
			groupID: "group-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `list users --group group-id`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("something"))
			},
			expErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			mc := &onepasswordclimock.OpCli{}
			test.mock(mc)

			repo, err := onepasswordcli.NewRepository(mc)
			require.NoError(err)

			gotMembership, err := repo.GetMembershipByID(context.TODO(), test.groupID, test.userID)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expMembership, gotMembership)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryDeleteMembership(t *testing.T) {
	tests := map[string]struct {
		membership model.Membership
		mock       func(m *onepasswordclimock.OpCli)
		expErr     bool
	}{
		"Delete a membership correctly, should delete the membership.": {
			membership: model.Membership{UserID: "user-id", GroupID: "group-id"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `remove user user-id group-id`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)
			},
		},

		"Having an error while calling the op CLI, should fail.": {
			membership: model.Membership{UserID: "user-id", GroupID: "group-id"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `remove user user-id group-id`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("something"))
			},
			expErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			mc := &onepasswordclimock.OpCli{}
			test.mock(mc)

			repo, err := onepasswordcli.NewRepository(mc)
			require.NoError(err)

			err = repo.DeleteMembership(context.TODO(), test.membership)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			mc.AssertExpectations(t)
		})
	}
}
