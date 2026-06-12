.PHONY: all clean

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
