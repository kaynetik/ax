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

# CLI COMMANDS

archive:
	./ax -arc-in ../tmp_to_archive -arc-pass on -arc-out ../tmp_archive_out

extract:
	./ax -arc-extract ../tmp_archive_out -arc-pass on

enc:
	./ax -enc-in ../tmp_archive_out

dec:
	./ax -dec-in ../tmp_archive_out

push:
	./ax -git-repo $(REPO) \
	-arc-in ../tmp_to_archive -arc-pass on -arc-out ../tmp_archive_out  \
	-enc-in ../tmp_archive_out