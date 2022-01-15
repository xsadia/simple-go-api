all: test vet build

test:
	go test -v

vet:
	go vet ./...

build: 
	go build -o main.out main.go