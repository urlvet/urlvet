# Security

## Admin Authentication

Access to `/admin` is protected by a password you set. The password is never stored in plain text — only a secure hash of it is kept.

### Setup

**Step 1 — Create your password hash**

Go to **https://argon2.online** and fill in:

| Field | Value |
|---|---|
| Plain Text Input | your chosen password |
| Salt | any random string (e.g. `openssl rand -base64 12`) |
| Parallelism Factor | `4` |
| Memory Cost | `65536` |
| Iterations | `3` |
| Hash Length | `32` |

Copy the **Encoded Form** output (starts with `$argon2id$`) into `server/.env`, wrapped in **single quotes** (required — the `$` signs would otherwise be interpreted as shell variables):

```
ADMIN_PASSWORD_HASH='$argon2id$v=19$...'
```

**Step 2 — Generate a signing secret**

```bash
openssl rand -hex 32
```

```
ADMIN_JWT_SECRET=<output>
```

### Logging in

Visit `/admin`, enter your password, and you'll get a session valid for 24 hours. Sessions are stored in the browser tab — closing the tab signs you out.

Login allows 5 attempts a minute per IP, separate from the API's general limit of 20 a minute. After that, it returns `429` until the minute is up. Use a long password anyway: the limit slows guessing, it doesn't stop it.
