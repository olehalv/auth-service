# olehalv-auth-service
En enkel Go Web Applikasjon som bygger, serverer og verfiserer JWTs basert på email/password

## Quickstart
Kompiler applikasjonen med `go build` og kjør executable. Eller kjør den ukompilert med `go run .`

Web applikasjonen lytter på `localhost:{PORT}`

## .env
Applikasjonen krever følgene .env variables

- PORT
- ADMIN_EMAIL
- ADMIN_PASS
- JWT_SECRET
- JWT_ISSUER

## Endepunkter
Det finnes 2 api endepunkter for autoriserings og henting av brukerinformasjon. I tillegg er det 2 .html sider hvor man kan logge inn og se brukerinformasjon

### api

- `/api/auth` - Krever JSON-body med `{email: string; pass: string}`, returnerer `{returnUrl: string; token: string}`
- `/api/user` - Fungerer som et endepunkt for å hente brukerinformasjon samt autorisering. Krever Bearer Authorization Header med JWT. Returnerer `{email: string}`

### static

- `/` `/index.html` - Serverer en enkel login page
- `/user.html` - Viser informasjon returnert av `/api/user`