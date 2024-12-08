# Launch Service

This is a basic Go service that provides CRUD endpoints for managing launch data.

## Requirements

- Go 1.18+
- No external database required (in-memory storage)

## Build and Run

```bash
go mod tidy
go build -o launch-service ./cmd
./launch-service
