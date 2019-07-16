mockgen:
	go get github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen

	mockgen -source targets/targets.go -destination targets/targets_mock.go -package targets -self_package github.com/coveooss/credentials-sync/targets
	mockgen -source credentials/sources.go -destination credentials/sources_mock.go -package credentials -self_package github.com/coveooss/credentials-sync/credentials