# auth-service
En enkel Go Web Applikasjon som bygger, serverer og verfiserer JWTs basert på email/password hentet fra en PostgreSQL database

## Quickstart
Kompiler applikasjonen med `go build` og kjør executable. Eller kjør den ukompilert med `go run .`

Web applikasjonen lytter på `{HOST}:{PORT}`

## .env
Applikasjonen krever følgene .env variables

- HOST
- PORT
- JWT_SECRET
- JWT_ISSUER
- PSQL_URL
- MAX_HTTP_REQUESTS_PER_MINUTE

## DB migrering
Alt av DB migrering ligger under `migrations/`, men er ikke satt opp til å migrere automatisk, dette må gjøres manuelt med verktøy som flyway

Eksempel: `flyway migrate -url="jdbc:postgresql://localhost:5432/DB_NAME" -user="USER" -password="PASS" -locations="filesystem:./migrations"`

## api

- `POST: /api/auth` - Krever JSON-body med `{email: string; pass: string}`, returnerer `{returnUrl: string; token: string}`
- `GET: /api/user` - Fungerer som et endepunkt for å hente brukerinformasjon samt å autorisere forespørsler. Krever Bearer Authorization Header med JWT. Returnerer `{email: string}`
- `POST: /api/user` - Legge til ny bruker, krever samme JSON body som `/api/auth` endepunktet. Returnerer 201 hvis opprettet
