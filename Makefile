version = $(shell git log -n 1 --format=%h)
built_at = $(shell date +"%Y-%m-%dT%T%z")

ldflags = -X MODULE_NAME/types.Version=$(version) -X MODULE_NAME/types.BuiltAt=$(built_at)

build = go build -ldflags="$(ldflags)"

build:
	$(build) .

build-race:
	$(build) -race -o ./tmp/main .

dev:
	air $(SERVICE) -d

test:
	go test ./...
