SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=pets-service

ifndef VERSION
  VERSION := git-$(shell git describe --always)
endif

BUILD_TIME=`date +%FT%T%z`

LDFLAGS=-a -tags netgo -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -w" 

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build ${LDFLAGS} -o ${BINARY} main.go

build_arm: $(SOURCES)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build ${LDFLAGS} -o ${BINARY} main.go

docker: clean ${BINARY}
	docker build -t robertojrojas/pets-service .

docker_arm: clean build_arm
	docker build -t robertojrojas/pets-service-arm -f Dockerfile.arm .

docker_clean: clean
	docker rmi robertojrojas/pets-service 

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
