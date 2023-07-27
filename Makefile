BINARY_NAME=gogit


build:
	go build -o ${BINARY_NAME} main.go

run: build
	./${BINARY_NAME}

clean:
	rm -r .gogit
