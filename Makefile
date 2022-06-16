EXEC_DIR = execuatables/

.PHONY: all build clean format test vet
all: build
check: format vet test

build: clean format vet test
	@mkdir -p $(EXEC_DIR)
	
	@echo "[+] Building the Linux version"
	@go build -ldflags "-s -w" -o $(EXEC_DIR)imgedit cmd/main.go

	@echo "[+] Packaging the Linux version"
	@zip -j $(EXEC_DIR)imgedit_Linux.zip $(EXEC_DIR)imgedit > /dev/null

	@echo "[+] Removing the Linux binary"
	@rm $(EXEC_DIR)imgedit

	@echo
	@echo "[+] Building the Windows version"
	@env GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o $(EXEC_DIR)imgedit.exe cmd/main.go

	@echo "[+] Packaging the Windows version"
	@zip -j $(EXEC_DIR)imgedit_Windows.zip $(EXEC_DIR)imgedit.exe > /dev/null

	@echo "[+] Removing the Windows binary"
	@rm $(EXEC_DIR)imgedit.exe

	@echo
	@echo "[+] Building the MacOS version"
	@env GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o $(EXEC_DIR)imgedit cmd/main.go

	@echo "[+] Packaging the MacOS version"
	@zip -j $(EXEC_DIR)imgedit_MacOS.zip $(EXEC_DIR)imgedit > /dev/null

	@echo "[+] Removing the MacOS binary"
	@rm $(EXEC_DIR)imgedit

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