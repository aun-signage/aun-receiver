default: build

save:
		godep save .

build:
		go build .

install:
		go install .

run:
		go run main.go
