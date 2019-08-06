package sync

import (
	"testing"

	"github.com/coveooss/credentials-sync/credentials"
)

func TestSyncCredentials(t *testing.T) {
	cred1, cred2 := credentials.NewSecretText(), credentials.NewSecretText()
	cred1.ID = "test1"
	cred2.ID = "test2"

	config := &Configuration{StopOnError: true, TargetParallelism: 1}
	targetController, target := setTargetMock(t, config, "target", []string{"test1"}, false)
	sourceController, _ := setSourceMock(t, config, []credentials.Credentials{cred1, cred2})
	defer targetController.Finish()
	defer sourceController.Finish()

	// Asserts that UpdateCredentials is called with `test1` and `test2`
	target.EXPECT().UpdateCredentials(cred1).Times(1)
	target.EXPECT().UpdateCredentials(cred2).Times(1)

	config.Sync()
}

func TestSyncCredentialsAndDeleteUnsynced(t *testing.T) {
	cred1, cred2 := credentials.NewSecretText(), credentials.NewSecretText()
	cred1.ID = "test1"
	cred2.ID = "test2"

	config := &Configuration{StopOnError: true, TargetParallelism: 1}
	targetController, target := setTargetMock(t, config, "target", []string{"test3"}, true)
	sourceController, _ := setSourceMock(t, config, []credentials.Credentials{cred1, cred2})
	defer targetController.Finish()
	defer sourceController.Finish()

	// Asserts that UpdateCredentials is called with `test1` and `test2`
	target.EXPECT().UpdateCredentials(cred1).Times(1)
	target.EXPECT().UpdateCredentials(cred2).Times(1)

	// Asserts that DeleteCredentials is called with `unsynced`
	target.EXPECT().DeleteCredentials("test3").Times(1)

	config.Sync()
}
