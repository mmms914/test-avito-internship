package mocks_test

//go:generate go run github.com/vektra/mockery/v2@latest --name=Repository --filename=repository_mock.go --dir=../ --output=.  --outpkg=mocks_test
//go:generate go run github.com/vektra/mockery/v2@latest --name=RoomRepository --filename=room_repository_mock.go --dir=../ --output=.  --outpkg=mocks_test
