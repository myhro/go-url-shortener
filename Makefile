BINARY = go-url-shortener
COVER_FILE = coverage.out
COVER_REPORT = coverage.html

build:
	go build -v -o $(BINARY)

clean:
	rm -f $(BINARY) $(COVER_FILE) $(COVER_REPORT)

coverage:
	go tool cover -html $(COVER_FILE) -o $(COVER_REPORT)

lint:
	golint -set_exit_status ./...

test:
	GIN_MODE=test go test -coverprofile $(COVER_FILE)
