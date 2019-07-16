package sync

import (
	"testing"

	"github.com/coveooss/credentials-sync/credentials"
	"github.com/coveooss/credentials-sync/targets"
	"github.com/golang/mock/gomock"
)

func setSourceMock(t *testing.T, config *Configuration) *credentials.MockSource {
	ctrl := gomock.NewController(t)
	source := credentials.NewMockSource(ctrl)
	sourceCollection := credentials.NewMockSourceCollection(ctrl)
	sourceCollection.EXPECT().AllSources().Return([]credentials.Source{source}).AnyTimes()

	config.SetSources(sourceCollection)
	return source
}

func setTargetMock(t *testing.T, config *Configuration) *targets.MockTarget {
	ctrl := gomock.NewController(t)
	target := targets.NewMockTarget(ctrl)
	targetCollection := targets.NewMockTargetCollection(ctrl)
	targetCollection.EXPECT().AllTargets().Return([]targets.Target{target}).AnyTimes()

	config.SetTargets(targetCollection)
	return target
}

func TestDeleteListOfCredentials(t *testing.T) {
	config := &Configuration{}
	setSourceMock(t, config)
	target := setTargetMock(t, config)
	target.EXPECT().GetExistingCredentials().Return([]string{"test123"}).AnyTimes()

	target.EXPECT().DeleteCredentials("test123").Return(nil)
}
