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

func setTargetMock(t *testing.T, config *Configuration) (*gomock.Controller, *targets.MockTarget) {
	ctrl := gomock.NewController(t)
	target := targets.NewMockTarget(ctrl)
	targetCollection := targets.NewMockTargetCollection(ctrl)
	targetCollection.EXPECT().AllTargets().Return([]targets.Target{target}).AnyTimes()

	config.SetTargets(targetCollection)
	return ctrl, target
}

func TestDeleteListOfCredentials(t *testing.T) {
	config := &Configuration{
		CredentialsToDelete: []string{"test1", "test-not-exist"},
	}
	setSourceMock(t, config)
	targetController, target := setTargetMock(t, config)
	defer targetController.Finish()

	// The existing creds is only `test1`
	target.EXPECT().GetExistingCredentials().Return([]string{"test1"}).AnyTimes()
	// Does not assert on GetName
	target.EXPECT().GetName().Return("").AnyTimes()
	// Asserts that DeleteCredentials is called with `test1`
	target.EXPECT().DeleteCredentials("test1").Return(nil)

	config.DeleteListOfCredentials(target)

}
