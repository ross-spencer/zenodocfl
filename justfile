# CLI helpers.

# help
help:
 @echo "Command line helpers for this project.\n"
 @just -l

# Run go linting
linting:
 - goimports -w .
 - go fmt ./...
 - go vet ./...
 - staticcheck ./...

# Run pre-commit
all-checks:
  pre-commit run --all-files

# Setup linting
setup:
  go install golang.org/x/tools/cmd/godoc@latest
  go install golang.org/x/tools/cmd/goimports@latest
  go install honnef.co/go/tools/cmd/staticcheck@latest

# Fix imports
fix-imports:
  goimports -w .

# Run tests
test:
  go test ./...

# Run tests without cache
test-nocache:
  go test -count=1 ./...

# Docs
docs:
  godoc -http 0.0.0.0:8000

# Motet example
lister-example:
 ./lister/lister -search motetcycle -results 10

# Motet records
lister-records:
 ./lister/lister -search motetcycle -results 300 > records.jsonl
 cat records.jsonl | wc -l
 mv records.jsonl gather/

# Motet example
lister-example-check:
 ./lister/lister -search motetcycle -results 300 -checklist

# Reset crate output
reset-output:
 rm -r {{`pwd`}}/crater/output
 git checkout {{`pwd`}}/crater/output
