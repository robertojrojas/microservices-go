SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=dogs-service

ifndef VERSION
  VERSION := git-$(shell git describe --always)
endif

BUILD_TIME=`date +%FT%T%z`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" -a -tags netgo -ldflags '-w' 

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build ${LDFLAGS} -o ${BINARY} main.go

build_arm: $(SOURCES)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build ${LDFLAGS} -o ${BINARY} main.go

docker: clean ${BINARY}
	docker build -t robertojrojas/dogs-service .

docker_arm: clean build_arm
	docker build -t robertojrojas/dogs-service-arm -f Dockerfile.arm .

docker_clean: clean
	docker rmi robertojrojas/dogs-service 

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
