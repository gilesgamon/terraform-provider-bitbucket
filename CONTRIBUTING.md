# Contributing

Thanks for your interest in improving the Bitbucket Terraform provider!

## Requirements

- [Go](https://go.dev/doc/install) >= 1.23
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [golangci-lint](https://golangci-lint.run/) (matching the version pinned in
  [`.github/workflows/lint.yml`](.github/workflows/lint.yml))

## Building

```bash
go build ./...
```

## Formatting, vetting and linting

```bash
make fmt          # gofmt -w
make vet          # go vet
golangci-lint run # static analysis (see .golangci.yml)
```

CI runs `gofmt`, `go vet`, `go build`, `go test`, and `golangci-lint` on every
pull request; please make sure they pass locally first.

## Testing

```bash
# Unit tests (no credentials required)
go test ./...

# Acceptance tests (create/read/update/delete against a real Bitbucket
# workspace — these cost API calls and require credentials)
export TF_ACC=1
export BITBUCKET_USERNAME=...
export BITBUCKET_PASSWORD=...
go test ./bitbucket/ -v -timeout 120m
```

Prefer credential-free unit tests for pure logic (flatten helpers, ID parsers,
client helpers). Reserve `TestAcc*` tests for behaviour that must hit the API.

## Documentation

- Reference docs live in [`docs/`](docs/), one page per resource/data source.
- A schema-driven scaffolder is available for new endpoints:

  ```bash
  go run ./tools/docgen
  ```

  It only creates pages that are missing, so existing hand-written docs are
  never overwritten. Enrich generated pages with descriptions and examples.
- `docs/**` changes are validated by markdownlint and markdown-link-check in CI.

## Adding a new endpoint

1. Add the resource/data source under `bitbucket/` following the existing
   conventions (schema, `flatten*`, ID helpers).
2. Register it in [`bitbucket/provider.go`](bitbucket/provider.go).
3. Add a docs page (or run `go run ./tools/docgen`).
4. Add an example under `examples/` where useful.
5. Update [`CHANGELOG.md`](CHANGELOG.md).

Only add endpoints that exist in the Bitbucket Cloud API
([`reference/swagger.v3.json`](reference/swagger.v3.json)).

## Pull requests

- Keep PRs focused; one logical change per commit.
- Use clear, descriptive commit messages.
- Do not force-push shared branches or amend others' commits.
