.PHONY: test
test:
	GIN_MODE=release go test `go list ./... | grep -vE '(/example)'`


.PHONY: cover
cover:
	GIN_MODE=release go test -cover `go list ./... | grep -vE '(/example)'`

.PHONY: coveralls
coveralls:
	GIN_MODE=release go test -covermode=atomic -coverprofile=coverage.out `go list ./... | grep -vE '(/example)'`

.PHONY: coverhtml
coverhtml:
	GIN_MODE=release go test -coverprofile=coverage.out `go list ./... | grep -vE '(/example)'`
	go tool cover -html=coverage.out