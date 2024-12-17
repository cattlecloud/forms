set shell := ["bash", "-u", "-c"]

export scripts := ".github/workflows/scripts"
export GOBIN := `echo $PWD/.bin`

# show available commands
[private]
default:
    @just --list

# tidy up Go modules
[group('build')]
tidy:
    go mod tidy

# run tests across source tree
[group('build')]
test:
    go test -v -race -count=1 ./...

# ensure copywrite headers present on source files
[group('lint')]
copywrite:
    copywrite \
        --config {{scripts}}/copywrite.hcl headers \
        --spdx "BSD-3-Clause"

# apply go vet command on source tree
[group('lint')]
vet:
    go vet ./...

# apply golangci-lint linters on source tree
[group('lint')]
lint: vet
    golangci-lint run --config .github/workflows/scripts/golangci.yaml

