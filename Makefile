BINARY_NAME=netbox-prefix-creator

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}_darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}_linux main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}_windows.exe main.go

clean:
	go clean
	rm ${BINARY_NAME}_darwin
	rm ${BINARY_NAME}_linux
	rm ${BINARY_NAME}_windows.exe
