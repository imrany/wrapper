all: build run

delete: 
	rm -rf bin
	rm -rf dist
	rm -rf ./wrapper

build:
	$(MAKE) delete
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/wrapper-linux main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/wrapper-windows.exe main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/wrapper-darwin main.go	

run:
	./wrapper --port=8080 --gemini-api-key=$(GEMINI_API_KEY)

dev: 
	CompileDaemon -build="go build -o wrapper main.go" -command="./wrapper --port=8080 --gemini-api-key=$(GEMINI_API_KEY)"

ensure-compile-daemon:
	@which go > /dev/null || (echo "Error: Go is not installed or not in PATH" && exit 1)
	@which CompileDaemon > /dev/null || (echo "Installing CompileDaemon..." && go install github.com/githubnemo/CompileDaemon@latest)

proto:
	cd proto && bash generate_proto.sh