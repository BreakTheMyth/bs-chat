.PHONY: all clean linux

ifdef OS
RM = rmdir /S /Q
else
RM = rm -rf
endif

all:
	cd web && npm run build
	go build

clean:
	go clean
	cd web && $(RM) build

linux:
	cd web && npm run build
	set GOOS=linux&& set GOARCH=amd64&& go build
