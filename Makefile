bin: bin/es_darwin bin/es_linux bin/es_windows

bin/es_darwin:
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -o bin/es_darwin cmd/es/*.go
	openssl sha512 bin/es_darwin > bin/es_darwin.sha512

bin/es_linux:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/es_linux cmd/es/*.go
	openssl sha512 bin/es_linux > bin/es_linux.sha512

bin/es_windows:
	mkdir -p bin
	GOOS=windows GOARCH=amd64 go build -o bin/es_windows cmd/es/*.go
	openssl sha512 bin/es_windows > bin/es_windows.sha512