install:
	@go mod download
	@go build -o $$GOPATH/bin/imgcov ./cmd/main.go

uninstall:
	@rm -f $$GOPATH/bin/imgcov