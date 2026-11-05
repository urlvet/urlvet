# Security Policy

## Reporting a Vulnerability

Please report security issues privately. **Don't open a public issue or PR.**

- **Preferred:** [Report a vulnerability](https://github.com/urlvet/urlvet/security/advisories/new) on GitHub. Only maintainers can see it.
- **Or email** [hi@url.vet](mailto:hi@url.vet) with "Security" in the subject. If you don't hear back within 3 days, check your spam folder and email again, or use the GitHub form.

Include:

- What the issue is and what an attacker could do with it
- Steps to reproduce, and a proof of concept if you have one
- The version or commit you tested, and how it was deployed (Docker, behind a proxy or CDN, etc.)

## What happens next

1. We reply within 3 days to say we've got it.
2. We check the report and tell you whether we can reproduce it.
3. We fix it and keep you updated. We aim to release a fix within 30 days of confirming the issue.
4. For confirmed issues, we publish a [GitHub Security Advisory](https://github.com/urlvet/urlvet/security/advisories) and request a CVE when it's warranted. We credit you unless you'd rather not be named.

Please keep the details private until the advisory is published. After that, you're welcome to publish your own write-up.

## Scope

In scope: the backend API, the website, the browser extension, and the Docker setup in this repository.

Not in scope:

- Values the website shows on purpose, such as `PUBLIC_*` settings in `/_app/env.js`
- Problems that need a misconfigured deployment, unless our docs told you to configure it that way
- Volumetric DoS, and reports from automated scanners without a working proof of concept

## Running your own instance

- Use the latest version from the `main` branch.
- Keep the backend behind a reverse proxy, and don't expose ports `8080`, `6379` or `9222` to the internet. See [docs/deployment.md](../docs/deployment.md).
- If a CDN or a proxy with a public IP connects straight to the backend, set `TRUSTED_PROXIES` or `TRUSTED_PLATFORM`. See [docs/configuration.md](../docs/configuration.md#client-ip-and-rate-limiting).
- If you use `TRUSTED_PLATFORM`, make sure only that platform can reach the backend, for example with a firewall that allows only its IP ranges.
