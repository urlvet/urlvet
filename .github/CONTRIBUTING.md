# Contributing to url.vet

Thank you for your interest in contributing!
We welcome bug reports, feature requests, code, documentation, and testing help.


## Project Structure

- `server/` – Backend (Golang)
- `web/website` – Frontend (Svelte)
- `docker/` – Dockerfiles & Compose configs
- `docs/` – Project documentation


## Contributor License Agreement

Before your first pull request is merged, you must agree to the [Contributor License Agreement](CLA.md). By opening a PR you confirm you have read and accepted it.


## How to Contribute

1. Fork the repo
2. Clone your fork and create a branch:
   ```bash
   git checkout -b your-feature-name
   ```
3. Make your changes
4. Run CI checks locally:
   ```bash
   make ci
   ```
5. Push to your fork and open a Pull Request

### PR checklist

- [ ] `make ci` passes locally (lint, format, type-check, tests)
- [ ] New behaviour is covered by tests where applicable
- [ ] PR description explains **what** changed and **why**
- [ ] No unrelated changes bundled in the same PR
- [ ] You have read and agree to the [CLA](CLA.md)


## Development Setup

Run `make help` to see all available commands.

### Check prerequisites

```bash
make doctor
```

### Dev environment (Docker — recommended)

```bash
make dev-build   # build dev images
make dev-up      # start the dev stack
make dev-logs    # tail logs
make dev-down    # stop
```

### Local (no Docker)

```bash
make tidy                 # go mod tidy
make local-run-backend    # run backend via Air (hot-reload)
make local-run-frontend   # run Svelte dev server
```


## Code Style & Tools

```bash
make format          # go fmt (backend)
make format-frontend # prettier (frontend)
make lint            # go vet
make check-frontend  # svelte-check + TypeScript
make ci              # runs all of the above + tests (backend + frontend)
```

Commit message convention:

```
feat(auth): add token expiration check
fix(api): correct 404 response logic
```


## Running Tests

```bash
make test            # backend (Go)
make test-frontend   # frontend (Vitest)
make ci              # both + lint + format + type-check
```

Frontend tests live in `web/website/src/` as `*.test.ts` files.
Run `make test-frontend` with `--ui` via `npm run test:ui` in `web/website/` for an interactive browser UI.


## Versioning

url.vet is in beta, so versions stay below 1.0 (`0.MINOR.PATCH`). After `0.9` comes `0.10`, not `1.0`.

- **Patch** (`0.4.1`): bug fixes, detection tweaks, list updates, small features. Upgrading needs no changes from self-hosters.
- **Minor** (`0.5.0`): released on a schedule, roughly every 3–6 months, bundling the work on `main` since the last one. Also used early when upgrading needs action from self-hosters, such as a changed default or a removed setting. The release notes say what to change.
- **Pre-releases** (`0.5.0-beta.1`, `0.5.0-rc.1`): for testing a minor release before it ships.
- **1.0** comes when the API response format, configuration and upgrade path are stable enough to promise not to break them. It isn't reached by counting.

The hosted site at url.vet deploys from `main` and doesn't wait for releases.

To release: update [CHANGELOG.md](../CHANGELOG.md), merge `dev` into `main`, tag the merge commit (`git tag -a v0.4.0 -m "v0.4.0 (beta)"`), push the tag, and create a GitHub Release marked as a pre-release.


## Reporting Bugs

1. Search existing issues first
2. If not found, open a new issue with:
   - Steps to reproduce
   - Logs or screenshots
   - Your environment (OS, browser, etc.)


<!-- ## Feature Requests

Use the [Feature Request template](../../issues/new?template=feature_request.yml)
Describe the use case and expected behavior. -->


## Code of Conduct

Please follow our [Code of Conduct](CODE_OF_CONDUCT.md).


## Thanks

Your contributions make this project better — whether it's code, feedback, or documentation.
