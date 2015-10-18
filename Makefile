PATH := ${PATH}:${GOPATH}/bin

all: deps build

install: deps
	go install
build:
	go build ./... 
	go build
imports: tools
	goimports -w .
	goimports -w */*.go
test:
	go test ./...
deps:
	go get -d -t
update-tools:
	go get -u github.com/golang/lint/golint
	go get -u golang.org/x/tools/cmd/goimports
tools:
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
clean:
	rm -rf reports/
cover: report-dir
	go test -coverprofile=reports/coverage.raw -covermode count 
bench: 
	go test -bench=.
check: report-dir tools
	golint ./... > reports/lint.txt
	go vet 2> reports/vet.txt
report-dir:
	mkdir -p reports

