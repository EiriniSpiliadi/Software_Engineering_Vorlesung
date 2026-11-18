# Ski-Verleih (Prototyp)
Nur Prototyp. Priorität: schnell verständlich, leicht änderbar. Entwicklung nach SDD: zuerst SPEC.md, dann Code.

## Stack
- Backend: Go (Structs + Methoden, Interfaces nur wo sinnvoll), DB: In-Memory (Map + Mutex)
- Frontend: reines HTML/JS (`web/index.html`), nur Calls gegen die Go-API

## Commands
- Start: `go run .` (http://localhost:8080)
- Tests: `go test ./...`
- Format: `go fmt ./...`

## Architektur
- `main.go` → Einstieg, Server, Demo-Daten
- `internal/domain/` → Structs + Methoden (Item, Person, Rental)
- `internal/store/` → thread-safe In-Memory-Store
- `internal/handler/` → HTTP-Handler (JSON-API unter `/api`)

## Konventionen
- Fehler mit `fmt.Errorf("...: %w", err)` wrappen; keine neuen Dependencies ohne Nachfrage
- Store muss thread-safe sein; keine Secrets im Code
- Commit-Messages: kurz und imperativ; bei Unklarheit fragen statt raten
