# Release Checklist

## Before tagging

- Run `make fmt`
- Run `make test`
- Run `make build`
- Run CLI smoke test: `scripts/smoke-test.sh`
- Verify `whisprgo version` contains expected version/commit/date
- Update `CHANGELOG.md`
- Ensure `VERSION` matches planned tag (for example `0.1.0` -> `v0.1.0`)

## Create release tag

- `git tag v$(cat VERSION)`
- `git push origin v$(cat VERSION)`

## Validate GitHub release artifacts

- Verify linux amd64 artifact exists
- Verify linux arm64 artifact exists
- Download and run `whisprgo version` from both artifacts
