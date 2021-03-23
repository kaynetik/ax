lint:
	gofumpt -w -s ./..
	golint ./...
	golangci-lint run --fix

test:
	go test ./...

build:
	go build -o ax cmd/cli/main.go

update_cache:
	curl https://sum.golang.org/lookup/github.com/kaynetik/ax@v$(VER)

# TEMP COMMANDS

tmp-archive:
	./ax -arc-in ../tmp_to_archive -arc-pass on -arc-out ../tmp_archive_out