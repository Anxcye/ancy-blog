# Origin Certificates

Place Cloudflare Origin Certificate files in this directory.
These files are mounted into the Caddy container at `/etc/caddy/certs`.

Recommended filenames:
- `origin.pem`
- `origin.key`

The certificate should cover:
- the apex domain (for example `example.com`)
- the wildcard subdomain (for example `*.example.com`)

After replacing these files, run:
```bash
cd deploy
./rotate-origin-cert.sh
```
