export GOPATH = $(shell pwd)

# GIT_VERSION is the only buildtime
all: 
all: 
	@go build -o myserver -ldflags "-X main.git_VERSION=$(shell git rev-parse HEAD)"

# I think this target shuold be a build dep of all, all: dependencies
# But I don't want it to build everytime, only if it's changed
# No local changes anyway....
dependencies:
	go get -u gopkg.in/dgrijalva/jwt-go.v3
	go get -u github.com/nu7hatch/gouuid
	go get -u github.com/autopogo/justClaims
	go get -u github.com/autopogo/justLogging

 #-ldflags "-X justClaims.git_VERSION=$(shell git rev-parse HEAD)"
