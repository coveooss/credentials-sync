package sync

import (
	"fmt"
	"testing"

	"github.com/coveooss/credentials-sync/credentials"
	"github.com/coveooss/credentials-sync/targets"
	"github.com/golang/mock/gomock"
)

func setSourceMock(t *testing.T, config *Configuration, creds []credentials.Credentials) (*gomock.Controller, *credentials.MockSourceCollection) {
	ctrl := gomock.NewController(t)
	source := credentials.NewMockSource(ctrl)
	sourceCollection := credentials.NewMockSourceCollection(ctrl)
	sourceCollection.EXPECT().AllSources().Return([]credentials.Source{source}).AnyTimes()

	if creds != nil {
		sourceCollection.EXPECT().Credentials().Return(creds, nil).AnyTimes()
	}

	config.SetSources(sourceCollection)
	return ctrl, sourceCollection
}

func setMultipleTargetMock(t *testing.T, config *Configuration, name string, credentials []string, shouldDeleteUnsynced bool, targetNumber int) (*gomock.Controller, []*targets.MockTarget) {
	ctrl := gomock.NewController(t)
	targetCollection := targets.NewMockTargetCollection(ctrl)

	targetsToReturn := []*targets.MockTarget{}
	targetsToReturnInterface := []targets.Target{}
	for i := 0; i < targetNumber; i++ {
		target := targets.NewMockTarget(ctrl)
		target.EXPECT().GetExistingCredentials().Return(credentials).AnyTimes()
		target.EXPECT().GetName().Return(fmt.Sprintf("%s-%v", name, i)).AnyTimes()
		target.EXPECT().GetTags().Return(map[string]string{}).AnyTimes()
		target.EXPECT().ToString().Return(fmt.Sprintf("%s-%v", name, i)).AnyTimes()
		target.EXPECT().ShouldDeleteUnsynced().Return(shouldDeleteUnsynced).AnyTimes()
		targetsToReturn = append(targetsToReturn, target)
		targetsToReturnInterface = append(targetsToReturnInterface, target)
	}
	targetCollection.EXPECT().AllTargets().Return(targetsToReturnInterface).AnyTimes()

	config.SetTargets(targetCollection)
	return ctrl, targetsToReturn
}

func setTargetMock(t *testing.T, config *Configuration, name string, credentials []string, shouldDeleteUnsynced bool) (*gomock.Controller, *targets.MockTarget) {
	ctrl, targetsToReturn := setMultipleTargetMock(t, config, name, credentials, shouldDeleteUnsynced, 1)
	targetsToReturn[0].EXPECT().Initialize(gomock.Any()).AnyTimes()
	return ctrl, targetsToReturn[0]
}
