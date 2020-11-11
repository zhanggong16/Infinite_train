ROOT_PATH=$(shell pwd)
PACKAGES := $$(go list ./...| grep -vE 'vendor')
FILES := $$(find . -name '*.go' | grep -vE 'vendor'| grep -vE 'Godeps'| grep -vE "jss" |grep -vE "mock")
VERSIONS = $(shell ./genver.sh show)


# target vm operating sysgem: linux/windows
export VMOSTYPE=linux


.PHONY:
all: check build


release: check fmt build_release


errcheck:
	go get github.com/kisielk/errcheck
	errcheck -ignoretests -verbose -blank $(PACKAGES)


check:
    # detect deadlock
	@ go vet ./...


fmt:
	@echo "gofmt (simplify)"
	@ gofmt -s -l -w $(FILES) 2>&1 | awk '{print} END{if(NR>0) {exit 1}}'


#build: version
build:
	export GOPROXY=https://goproxy.cn,direct && \
	export GOBIN=${ROOT_PATH}/bin && \
	go install  -gcflags "-N -l"  ./...

build_release: version build


unittest:
	go tool cover -func=total.cov
	go tool cover -html=total.cov -o=coverage-report.html


.PHONY:
test:
	export GOBIN=${ROOT_PATH}/bin && \
	go test --race -timeout 2m $(PACKAGES)


.PHONY:
clean:
	@rm -rf ./output
	@rm -rf ./bin/*
