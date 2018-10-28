export GOPATH = $(shell pwd)

all: export GIT_VERSION=$(git rev-parse HEAD) # fuck a tag
all: export COOKIE_KEY ?= zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz
all: export JWT_KEY ?= xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
all:
		@echo "Did your environment set the cookie and jwt keys?"
		@go build -o myserver -ldflags "-X main.git_VERSION=$(GIT_VERSION) -X main.key_COOKIE_SIGNING=$(COOKIE_KEY) -X main.key_JWT_SIGNING=$(JWT_KEY)"

dependencies:
		go get -u github.com/gorilla/sessions
		go get -u github.com/gorilla/context
		go get -u gopkg.in/dgrijalva/jwt-go.v3
		go get -u github.com/nu7hatch/gouuid
