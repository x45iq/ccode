build:
	mkdir -p bin
	go build -o bin/ccode .

lint:
	golangci-lint run ./...

test:
	go test ./... -v

clean:
	rm -rf bin

