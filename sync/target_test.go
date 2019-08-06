package sync

import (
	"testing"

	"github.com/coveooss/credentials-sync/credentials"
	"github.com/golang/mock/gomock"
)

func TestDeleteListOfCredentials(t *testing.T) {
	config := &Configuration{
		CredentialsToDelete: []string{"test1", "test-not-exist"},
	}
	targetController, target := setTargetMock(t, config, "", []string{"test1"}, false)
	defer targetController.Finish()

	// Asserts that DeleteCredentials is called with `test1`
	target.EXPECT().DeleteCredentials("test1").Return(nil)

	config.DeleteListOfCredentials(target)
}

func TestUpdateListOfCredentials(t *testing.T) {
	config := NewConfiguration()
	targetController, target := setTargetMock(t, config, "", []string{"test1"}, false)
	defer targetController.Finish()

	cred1, cred2 := credentials.NewSecretText(), credentials.NewSecretText()
	cred1.ID = "test1"
	cred2.ID = "test2"

	// Asserts that UpdateCredentials is called with `test1` and `test2`
	target.EXPECT().UpdateCredentials(cred1).Times(1)
	target.EXPECT().UpdateCredentials(cred2).Times(1)

	config.UpdateListOfCredentials(target, []credentials.Credentials{cred1, cred2})
}

func TestDeleteUnsyncedCredentials(t *testing.T) {
	config := NewConfiguration()
	targetController, target := setTargetMock(t, config, "", []string{"unsynced"}, true)
	defer targetController.Finish()

	cred1, cred2 := credentials.NewSecretText(), credentials.NewSecretText()
	cred1.ID = "test1"
	cred2.ID = "test2"

	// Do not verify this call
	target.EXPECT().UpdateCredentials(gomock.Any()).AnyTimes()

	// Asserts that DeleteCredentials is called with `unsynced`
	target.EXPECT().DeleteCredentials("unsynced").Times(1)

	config.UpdateListOfCredentials(target, []credentials.Credentials{cred1, cred2})
}
