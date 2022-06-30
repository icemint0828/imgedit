EXEC_DIR = execuatables/

.PHONY: all build dev check clean format test vet
all: build
check: format vet test

dev: check
	@go build -o $$GOPATH/bin/imgedit ./cmd/main.go

build: clean check
	@mkdir -p $(EXEC_DIR)

	@echo "[+] Copy license"
	@cp LICENSE $(EXEC_DIR)LICENSE
	@cp THIRD_PARTY_LICENCES $(EXEC_DIR)THIRD_PARTY_LICENCES
	
	@echo "[+] Building the Linux version"
	@go build -ldflags "-s -w" -o $(EXEC_DIR)imgedit cmd/main.go

	@echo "[+] Packaging the Linux version"
	@zip -j $(EXEC_DIR)imgedit_Linux.zip -r $(EXEC_DIR)imgedit $(EXEC_DIR)LICENSE $(EXEC_DIR)THIRD_PARTY_LICENCES > /dev/null

	@echo "[+] Removing the Linux binary"
	@rm $(EXEC_DIR)imgedit

	@echo
	@echo "[+] Building the Windows version"
	@env GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o $(EXEC_DIR)imgedit.exe cmd/main.go

	@echo "[+] Packaging the Windows version"
	@zip -j $(EXEC_DIR)imgedit_Windows.zip -r $(EXEC_DIR)imgedit.exe $(EXEC_DIR)LICENSE $(EXEC_DIR)THIRD_PARTY_LICENCES > /dev/null

	@echo "[+] Removing the Windows binary"
	@rm $(EXEC_DIR)imgedit.exe

	@echo
	@echo "[+] Building the MacOS version"
	@env GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o $(EXEC_DIR)imgedit cmd/main.go

	@echo "[+] Packaging the MacOS version"
	@zip -j $(EXEC_DIR)imgedit_MacOS.zip -r $(EXEC_DIR)imgedit $(EXEC_DIR)LICENSE $(EXEC_DIR)THIRD_PARTY_LICENCES > /dev/null

	@echo "[+] Removing the MacOS binary"
	@rm $(EXEC_DIR)imgedit

	@echo "[+] Removing license"
	@rm $(EXEC_DIR)LICENSE $(EXEC_DIR)THIRD_PARTY_LICENCES

	@echo "[+] Done"

clean:
	@echo "[+] Cleaning files"
	@rm -rf $(EXEC_DIR)
	@echo "[+] Done"
	@echo

format:
	@echo "[+] Formatting files"
	@gofmt -w *.go

vet:
	@echo "[+] Running Go vet"
	@go vet

test:
	@echo "[+] Running tests"
	@go test

gitHubActions:
	@echo "[+] Building container image - GitHub Actions"
	@env GOOS=linux CGO_ENABLED=0 go build --ldflags '-s -w' -o imgedit cmd/main.go && chmod +x imgedit
	@echo "[+] Done"
