mockgen:
	go get github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen
	
	rm targets/target_mock_test.go
	mockgen -source targets/targets.go -destination targets/target_mock_test.go -package targets