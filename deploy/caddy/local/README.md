# Local Caddy Overrides

Put server-specific Caddy snippets in this directory.
These files are mounted into the Caddy container at `/etc/caddy/local/*.caddy` and are imported by the main `Caddyfile`.

Examples:
- exact legacy redirects
- temporary maintenance redirects
- special headers for one environment

Recommended file naming:
- `legacy-redirects.caddy`
- `maintenance.caddy`
- `admin-access.caddy`
