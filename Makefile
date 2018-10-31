export GOPATH = $(shell pwd)

# GIT_VERSION is the only buildtime
all: export GIT_VERSION=$(git rev-parse HEAD) # fuck a tag
all: 
		@go build -o myserver -ldflags "-X main.git_VERSION=$(GIT_VERSION)"

# I think this target shuold be a build dep of all, all: dependencies
# But I don't want it to build everytime, only if it's changed
# No local changes anyway....
dependencies:
		go get -u gopkg.in/dgrijalva/jwt-go.v3
		go get -u github.com/nu7hatch/gouuid
		go get -u github.com/ayjayt/justTheClaims
