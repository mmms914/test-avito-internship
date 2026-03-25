package mocks_test

//go:generate go run github.com/vektra/mockery/v2@latest --name=Repository --filename=repository_mock.go --dir=../ --output=. --outpkg=mocks_test
//go:generate go run github.com/vektra/mockery/v2@latest --name=SlotRepository --filename=slot_repository_mock.go --dir=../ --output=. --outpkg=mocks_test
//go:generate go run github.com/vektra/mockery/v2@latest --name=LinkManager --filename=link_manager_mock.go --dir=../ --output=. --outpkg=mocks_test
//go:generate go run github.com/vektra/mockery/v2@latest --name=Logger --filename=logger_mock.go --dir=../ --output=. --outpkg=mocks_test
