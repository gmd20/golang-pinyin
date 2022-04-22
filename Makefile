export PATH := /home/ming/go/bin:/home/ming/go_project/bin:$(PATH)
export GOROOT=/home/ming/go
export GOPATH=/home/ming/go_project
export GOPROXY=https://goproxy.cn,direct
export GO111MODULE=on
# LDFLAGS := -s -w # https://docs.studygolang.com/cmd/link/

all:
	go build -trimpath -ldflags "$(LDFLAGS)" -o ./bin/pinyin ./cmd/pinyin
	
fmt:
	go fmt ./...


test:
	go test -v --cover ./cmd/...

clean:
	rm -f ./bin/pinyin
