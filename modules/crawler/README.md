# crawler

Small, experimental Go web crawler — quick and low-ceremony.

Quick start

```bash
# needs Go (1.20+)
go run main.go <BASE_URL> <MAX_CONCURRENT> <MAX_PAGE>

# build and run
go build -o crawler .
./crawler <BASE_URL> <MAX_CONCURRENT> <MAX_PAGE>

# run tests
go test -v ./...
```

What it does

- Normalizes URLs (`normalize_url.go`).
- Parses page data (`parser.go`, `page_data.go`).
- Small test suite in `*_test.go` files.

