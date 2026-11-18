# Spezifikation Fahrrad-Verleih
## Funktional
- F1 Dashboard: verfügbare / ausgeliehene / gesamte Fahrräder, registrierte Personen; Schnellzugriff (verwalten, Personen, neue Ausleihe)
- F2 Fahrrad-Verwaltung: hinzufügen (Nr. + Typ, Nr. eindeutig), löschen (nicht während Ausleihe)
- F3 Personen: registrieren, auflisten
- F4 Ausleihe: Fahrrad + Person, frei wählbarer Beginn, Dauer in Stunden, Rückgabe
- F5 Preise: Stundenpreis pro Typ änderbar; pro Ausleihe überschreibbar; Gesamt = Stunden × Preis/h
## Nicht-funktional
- In-Memory, thread-safe; benutzerfreundliche, einfache Oberfläche; Fehlermeldungen verständlich
## Akzeptanz
- Doppelte Ausleihe desselben Artikels wird abgelehnt; `go test ./...` grün

## Usability requirements
- U1 Every action gives visible feedback (success or error message)
- U2 Destructive actions (delete, return) require confirmation
- U3 Rental form shows end time and total price before submitting
- U4 Lists are searchable; empty lists show an explanation
- U5 Status is shown by text and color (available, rented, overdue)
- U6 The UI works on small screens and with keyboard only

## Business rules (rentals)
- BR1 Item number is required and unique; item type must be a configured type
- BR2 A rented item cannot be rented again or deleted
- BR3 Person name is required and unique (case-insensitive)
- BR4 Duration: 1 to 72 hours; start at most 7 days in the future
- BR5 One type at a time: while a person has an active rental, they can only rent items of the same type until everything is returned
- BR6 Maximum 4 active items per person
- BR7 A person with an overdue rental (end time passed, not returned) cannot rent anything new
- BR8 Price per hour defaults to the type's price; it can be overridden per rental
- BR9 Returning makes the item available again
Limits (72 h, 4 items) are constants in `internal/domain/domain.go`.

## Age rule
- BR10 Persons have an age (1-120). A person under 18 needs a supervisor for every rental: a different, registered person aged 18 or older
- The supervisor is stored with the rental (`supervisor_id`); the UI shows the supervisor field only for minors

## Pricing rules
- P1 Each item type has a price per hour (default shown on the Preise page); it can be overridden per rental
- P2 Duration discount: 10 % from 4 h, 20 % from 8 h (constants `Tiers` in `internal/domain/domain.go`)
- P3 Total = hours x price per hour x (1 - discount), rounded to cents
- P4 The rental form shows end time and the price calculation before confirming; the Preise page shows example prices for 1, 2, 4, 8 and 24 hours
- Example (5 EUR/h): 2 h = 10.00 EUR, 4 h = 18.00 EUR, 8 h = 32.00 EUR
