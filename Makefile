COVER_FILE = coverage.out
COVER_REPORT = coverage.html

coverage:
	go tool cover -html $(COVER_FILE) -o $(COVER_REPORT)

test:
	GIN_MODE=test go test -coverprofile $(COVER_FILE)
