package sync

import (
	"fmt"
	"testing"

	"github.com/coveooss/credentials-sync/credentials"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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

	assert.Nil(t, config.Sync())
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

	assert.Nil(t, config.Sync())
}

func TestSyncCredentialsAndDeleteUnsyncedWithContinueOnError(t *testing.T) {
	cred1, cred2 := credentials.NewSecretText(), credentials.NewSecretText()
	cred1.ID = "test1"
	cred2.ID = "test2"

	config := &Configuration{StopOnError: false, TargetParallelism: 1}
	targetController, targets := setMultipleTargetMock(t, config, "target", []string{"test3", "test4"}, true, 2)
	sourceController, _ := setSourceMock(t, config, []credentials.Credentials{cred1, cred2})
	defer targetController.Finish()
	defer sourceController.Finish()

	// first target returns an error on initialize so it is removed from the pool
	targets[0].EXPECT().Initialize(gomock.Any()).Return(fmt.Errorf("Dummy error")).AnyTimes()
	targets[1].EXPECT().Initialize(gomock.Any()).Return(nil).AnyTimes()

	// Second target fails update on first cred but succeeds on second
	targets[1].EXPECT().UpdateCredentials(cred1).Return(fmt.Errorf("Dummy error")).Times(1)
	targets[1].EXPECT().UpdateCredentials(cred2).Return(nil).Times(1)

	// Second target fails delete on first cred but succeeds on second
	targets[1].EXPECT().DeleteCredentials("test3").Return(fmt.Errorf("Dummy error")).Times(1)
	targets[1].EXPECT().DeleteCredentials("test4").Return(nil).Times(1)

	assert.Nil(t, config.Sync())
}

func TestSyncCredentialsFailOnInitialize(t *testing.T) {
	cred1, cred2 := credentials.NewSecretText(), credentials.NewSecretText()
	cred1.ID = "test1"
	cred2.ID = "test2"

	config := &Configuration{StopOnError: true, TargetParallelism: 1}
	targetController, targets := setMultipleTargetMock(t, config, "target", []string{}, true, 2)
	sourceController, _ := setSourceMock(t, config, []credentials.Credentials{cred1, cred2})
	defer targetController.Finish()
	defer sourceController.Finish()

	targets[0].EXPECT().Initialize(gomock.Any()).Return(fmt.Errorf("Dummy error1")).AnyTimes()
	targets[1].EXPECT().Initialize(gomock.Any()).Return(nil).AnyTimes()

	assert.EqualError(t, config.Sync(), "Target `target` has failed initialization: Dummy error1")
}

func TestSyncCredentialsFailOnCredentialsUpdate(t *testing.T) {
	cred1, cred2 := credentials.NewSecretText(), credentials.NewSecretText()
	cred1.ID = "test1"
	cred2.ID = "test2"

	config := &Configuration{StopOnError: true, TargetParallelism: 1}
	targetController, targets := setMultipleTargetMock(t, config, "target", []string{}, true, 2)
	sourceController, _ := setSourceMock(t, config, []credentials.Credentials{cred1, cred2})
	defer targetController.Finish()
	defer sourceController.Finish()

	// Initialize works for both targets
	targets[0].EXPECT().Initialize(gomock.Any()).Return(nil).AnyTimes()
	targets[1].EXPECT().Initialize(gomock.Any()).Return(nil).AnyTimes()

	// First target succeeds but second target fails
	targets[0].EXPECT().UpdateCredentials(gomock.Any()).Return(nil).AnyTimes()
	targets[1].EXPECT().UpdateCredentials(cred1).Return(fmt.Errorf("Dummy error2")).Times(1)

	assert.EqualError(t, config.Sync(), "Failed to send credentials with ID test1 to target: Dummy error2")
}

func TestSyncCredentialsFailOnCredentialsDelete(t *testing.T) {
	cred1, cred2 := credentials.NewSecretText(), credentials.NewSecretText()
	cred1.ID = "test1"
	cred2.ID = "test2"

	config := &Configuration{StopOnError: true, TargetParallelism: 1}
	targetController, targets := setMultipleTargetMock(t, config, "target", []string{"test3", "test4"}, true, 2)
	sourceController, _ := setSourceMock(t, config, []credentials.Credentials{cred1, cred2})
	defer targetController.Finish()
	defer sourceController.Finish()

	// Initialize works for both targets
	targets[0].EXPECT().Initialize(gomock.Any()).Return(nil).AnyTimes()
	targets[1].EXPECT().Initialize(gomock.Any()).Return(nil).AnyTimes()

	// UpdateCredentials works for both targets
	targets[0].EXPECT().UpdateCredentials(gomock.Any()).Return(nil).AnyTimes()
	targets[1].EXPECT().UpdateCredentials(gomock.Any()).Return(nil).AnyTimes()

	// Delete credentials fails for first target but succeeds for second
	targets[0].EXPECT().DeleteCredentials("test3").Return(fmt.Errorf("Dummy error3")).Times(1)
	targets[1].EXPECT().DeleteCredentials(gomock.Any()).Return(nil).AnyTimes()

	assert.EqualError(t, config.Sync(), "Failed to delete credentials with ID test3 to target: Dummy error3")
}
