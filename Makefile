BINARY_NAME=gogit


build:
	go build -o ./bin/${BINARY_NAME} main.go

run: build
	./bin/${BINARY_NAME}

clean:
	rm -r .gogit
