.DEFAULT_GOAL := mkvbot.exe
VERSION := 0.1.0

.PHONY: generate
generate:
	go generate ./...

.PHONY: mkvbot
mkvbot: generate
	GOOS=linux GOARCH=arm64 go build -o $@ -tags linux -ldflags "-X main.Version=$(VERSION)"

.PHONY: mkvbot.exe
mkvbot.exe: generate
	GOOS=windows go build -o $@ -tags windows -ldflags "-X main.Version=$(VERSION)"

.PHONY: clean
clean:
	rm -rf mkvbot.exe mkvbot dist
