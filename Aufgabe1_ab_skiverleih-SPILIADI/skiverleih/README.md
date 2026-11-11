# Ski-Verleih

A small rental system prototype for skis, built with **Spec-Driven Development (SDD)** for the Software Engineering course (Aufgabe 1). Go backend with an in-memory store and a plain HTML/JS dashboard.

## Features
- **Dashboard**: available, rented and total skis, registered persons, quick-access buttons
- **Ski management**: add by number + type (numbers are unique), delete (not while rented)
- **Persons**: register (name, age) and list
- **Rentals**: choose ski and person, free start time, duration in hours, return
- **Flexible pricing**: hourly price per type, editable; can be overridden per rental. Total = hours × price per hour

## Pricing
- Price per hour depends on the item type and can be changed on the "Preise" page; it can also be overridden per rental
- Duration discount: **10 % from 4 hours, 20 % from 8 hours** (applies to the whole rental)
- Total = hours × price per hour × (1 − discount), rounded to cents; the rental form shows the calculation before you confirm
- Example prices with the default hourly prices:

| Type | Price/h | 1 h | 2 h | 4 h | 8 h | 24 h |
|---|---|---|---|---|---|---|
| Alpin | 8.00 € | 8.00 € | 16.00 € | 28.80 € | 51.20 € | 153.60 € |
| Langlauf | 5.00 € | 5.00 € | 10.00 € | 18.00 € | 32.00 € | 96.00 € |
| Snowboard | 7.00 € | 7.00 € | 14.00 € | 25.20 € | 44.80 € | 134.40 € |

## Tech stack
| Part | Choice |
|---|---|
| Backend | Go 1.22+ (standard library only, no dependencies) |
| Storage | In-memory (`map` + `sync.RWMutex`), data is lost on restart |
| Frontend | Plain HTML/JS in `web/index.html`, talks only to the REST API |

## Getting started
```bash
go test ./...   # run tests
go run .        # start server
```
Open http://localhost:8080. The app starts with a few demo items and persons.

## Project structure
```
main.go                      entry point, server, demo data
internal/domain/             structs + methods (Item, Person, Rental)
internal/store/              thread-safe in-memory store (+ tests)
internal/handler/            HTTP handlers (JSON API under /api)
web/index.html               dashboard UI
CLAUDE.md                    project rules for the AI assistant
SPEC.md                      specification (written before the code)
```

## REST API
| Method | Path | Description |
|---|---|---|
| GET | `/api/dashboard` | counts for the dashboard |
| GET, POST | `/api/items` | list / add item (`number`, `type`) |
| DELETE | `/api/items/{id}` | delete item |
| GET, POST | `/api/persons` | list / add person (`name`, `age`) |
| GET, POST | `/api/rentals` | list / start rental (`item_id`, `person_id`, `start`, `hours`, `price_per_hour`, `supervisor_id` for persons under 18) |
| POST | `/api/rentals/{id}/return` | return item |
| GET | `/api/tiers` | duration discounts (min hours, discount) |
| GET | `/api/prices` | hourly price per type |
| PUT | `/api/prices/{type}` | set hourly price (`price`) |

Errors return HTTP 400 with a readable message (e.g. renting an item that is already rented).

## Spec-Driven Development
1. `SPEC.md` defines requirements and acceptance criteria.
2. `CLAUDE.md` gives the AI assistant the stack, commands, architecture and conventions.
3. Code is generated and refined against the spec; `go test ./...` checks the acceptance criteria.

## Limitations
- No persistence and no authentication (prototype)
- Prices are in EUR, no time zones beyond the browser's local time
- Rentals are not ended automatically; overdue rentals are flagged in the UI and block new rentals until returned

## Usability
- Dashboard cards are clickable shortcuts; hints guide first-time users (e.g. "create a person first")
- Active menu item is highlighted; Enter submits forms; first field is focused automatically
- Search by number or type; colored status badges (available / rented / active / overdue)
- "Rent" button directly in the item list pre-selects the item
- Live preview of end time and total price before confirming a rental
- Confirmation dialogs for deleting and returning; toast messages for success and errors
- Empty-state messages, responsive layout for small screens, visible keyboard focus

## Rental rules
- A person can only rent **one type at a time**: with an active rental, only items of the same type are offered until everything is returned
- Maximum **4 items** per person and **72 hours** per rental; start at most 7 days ahead
- A person with an **overdue** rental must return it before renting again
- Item numbers and person names must be unique; rented items cannot be deleted or rented twice
- A person under 18 can only rent with a registered adult (18+) as supervisor; the rental form then shows a supervisor field
- The rental form only offers allowed items and explains why others are hidden
