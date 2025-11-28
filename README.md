# Notes API

A simple REST API for managing notes. This is still a work in progress

## What it does

Right now you can:

-   Check health status (`GET /health`)
-   Create notes (`POST /notes`)
-   List all notes (`GET /notes`)
-   Get a note by ID (`GET /notes/:id`)

Still TODO:

-   Delete notes
-   Update notes
-   Proper persistence

## Running it

```bash
go mod download
go run cmd/service/main.go
```

The server will start on port 8080 by default, or whatever you set in the `PORT` environment variable.

## Example usage

Create a note:

```bash
curl -X POST http://localhost:8080/notes \
  -H "Content-Type: application/json" \
  -d '{"title": "My first note", "content": "This is some content"}'
```

Get all notes:

```bash
curl http://localhost:8080/notes
```

Get a specific note:

```bash
curl http://localhost:8080/notes/{note-id}
```

## Notes

-   Data is not persistent (MemoryStore)
-   Have nott decided on the final API structure, might change
