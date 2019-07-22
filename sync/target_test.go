package sync

import (
	"testing"

	"github.com/coveooss/credentials-sync/credentials"
	"github.com/coveooss/credentials-sync/targets"
	"github.com/golang/mock/gomock"
)

func setSourceMock(t *testing.T, config *Configuration) (*gomock.Controller, *credentials.MockSource) {
	ctrl := gomock.NewController(t)
	source := credentials.NewMockSource(ctrl)
	sourceCollection := credentials.NewMockSourceCollection(ctrl)
	sourceCollection.EXPECT().AllSources().Return([]credentials.Source{source}).AnyTimes()

	config.SetSources(sourceCollection)
	return ctrl, source
}

func setTargetMock(t *testing.T, config *Configuration, name string, credentials []string, shouldDeleteUnsynced bool) (*gomock.Controller, *targets.MockTarget) {
	ctrl := gomock.NewController(t)
	target := targets.NewMockTarget(ctrl)
	targetCollection := targets.NewMockTargetCollection(ctrl)
	targetCollection.EXPECT().AllTargets().Return([]targets.Target{target}).AnyTimes()

	target.EXPECT().GetExistingCredentials().Return(credentials).AnyTimes()
	target.EXPECT().GetName().Return(name).AnyTimes()
	target.EXPECT().ShouldDeleteUnsynced().Return(shouldDeleteUnsynced).AnyTimes()

	config.SetTargets(targetCollection)
	return ctrl, target
}

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
	config := &Configuration{}
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
	config := &Configuration{}
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
