# auth-service
En enkel Go Web Applikasjon som bygger, serverer og verfiserer JWTs basert på email/password hentet fra en PostgreSQL database.

## Quickstart
Kompiler applikasjonen med `go build` og kjør executable. Eller kjør den ukompilert med `go run .`

Web applikasjonen lytter på `{HOST}:{PORT}`

## .env
Applikasjonen krever følgene .env variables

- HOST (må defineres opp i /etc/hosts)
- PORT
- JWT_SECRET
- JWT_ISSUER
- PSQL_URL
- MAX_HTTP_REQUESTS_PER_MINUTE

## DB migrering
Alt av DB migrering ligger under `migrations/`, men er ikke satt opp til å migrere automatisk, dette må gjøres manuelt

## Endepunkter
Det finnes 2 api endepunkter for autoriserings og henting av brukerinformasjon. I tillegg er det 2 .html sider hvor man kan logge inn og se brukerinformasjon

### api

- `/api/auth` - Krever JSON-body med `{email: string; pass: string}`, returnerer `{returnUrl: string; token: string}`
- `/api/user` - Fungerer som et endepunkt for å hente brukerinformasjon samt autorisering. Krever Bearer Authorization Header med JWT. Returnerer `{email: string}`

### static

- `/` `/index.html` - Serverer en enkel login page
- `/user.html` - Viser informasjon returnert av `/api/user`
