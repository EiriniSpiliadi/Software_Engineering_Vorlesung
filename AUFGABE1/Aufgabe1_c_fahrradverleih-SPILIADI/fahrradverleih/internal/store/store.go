package store

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"fahrradverleih/internal/domain"
)

// Store is a thread-safe in-memory store.
type Store struct {
	mu      sync.RWMutex
	next    int
	items   map[int]*domain.Item
	persons map[int]*domain.Person
	rentals map[int]*domain.Rental
	prices  map[string]float64 // price per hour by item type
}

func New(types []string, price float64) *Store {
	s := &Store{next: 1, items: map[int]*domain.Item{}, persons: map[int]*domain.Person{},
		rentals: map[int]*domain.Rental{}, prices: map[string]float64{}}
	for _, t := range types {
		s.prices[t] = price
	}
	return s
}

func (s *Store) id() int { n := s.next; s.next++; return n } // call under lock

func bad(format string, a ...any) error {
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, a...), domain.ErrInvalid)
}

func (s *Store) AddItem(number, typ string) (domain.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if number == "" {
		return domain.Item{}, bad("Nr. fehlt")
	}
	if _, ok := s.prices[typ]; !ok {
		return domain.Item{}, bad("unbekannter Typ %q", typ)
	}
	for _, i := range s.items {
		if i.Number == number {
			return domain.Item{}, bad("Nr. %s existiert bereits", number)
		}
	}
	it := &domain.Item{ID: s.id(), Number: number, Type: typ}
	s.items[it.ID] = it
	return *it, nil
}

func (s *Store) DeleteItem(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	it, ok := s.items[id]
	if !ok {
		return bad("nicht gefunden")
	}
	if it.Rented {
		return bad("derzeit ausgeliehen")
	}
	delete(s.items, id)
	return nil
}

func (s *Store) Items() []domain.Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Item{}
	for _, i := range s.items {
		out = append(out, *i)
	}
	sort.Slice(out, func(a, b int) bool { return out[a].ID < out[b].ID })
	return out
}

func (s *Store) AddPerson(name string, age int) (domain.Person, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if name == "" {
		return domain.Person{}, bad("Name fehlt")
	}
	if age < 1 || age > 120 {
		return domain.Person{}, bad("Alter muss zwischen 1 und 120 liegen")
	}
	for _, x := range s.persons {
		if strings.EqualFold(x.Name, name) {
			return domain.Person{}, bad("Person %s ist bereits registriert", name)
		}
	}
	p := &domain.Person{ID: s.id(), Name: name, Age: age}
	s.persons[p.ID] = p
	return *p, nil
}

func (s *Store) Persons() []domain.Person {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Person{}
	for _, p := range s.persons {
		out = append(out, *p)
	}
	sort.Slice(out, func(a, b int) bool { return out[a].ID < out[b].ID })
	return out
}

// Rent starts a rental. price<=0 uses the default price of the item type; zero start means now.
func (s *Store) Rent(itemID, personID int, start time.Time, hours int, price float64, supervisorID int) (domain.Rental, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	it, ok := s.items[itemID]
	if !ok {
		return domain.Rental{}, bad("Artikel nicht gefunden")
	}
	if it.Rented {
		return domain.Rental{}, bad("Artikel ist bereits ausgeliehen")
	}
	if _, ok := s.persons[personID]; !ok {
		return domain.Rental{}, bad("Person nicht gefunden")
	}
	if hours < 1 || hours > domain.MaxHours {
		return domain.Rental{}, bad("Dauer muss zwischen 1 und %d Stunden liegen", domain.MaxHours)
	}
	now := time.Now()
	if start.After(now.Add(7 * 24 * time.Hour)) {
		return domain.Rental{}, bad("Beginn darf höchstens 7 Tage in der Zukunft liegen")
	}
	name := s.persons[personID].Name
	if s.persons[personID].Age < domain.AdultAge {
		sup, ok := s.persons[supervisorID]
		if !ok || supervisorID == personID || sup.Age < domain.AdultAge {
			return domain.Rental{}, bad("%s ist unter %d Jahre alt und braucht eine volljährige Aufsichtsperson", name, domain.AdultAge)
		}
	}
	active := 0
	for _, x := range s.rentals {
		if x.PersonID != personID || x.Returned {
			continue
		}
		active++
		if t := s.items[x.ItemID].Type; t != it.Type {
			return domain.Rental{}, bad("%s hat noch eine aktive Ausleihe vom Typ %s und kann erst nach der Rückgabe Typ %s ausleihen", name, t, it.Type)
		}
		if x.End().Before(now) {
			return domain.Rental{}, bad("%s hat eine überfällige Ausleihe – bitte zuerst zurückgeben", name)
		}
	}
	if active >= domain.MaxItemsPerPerson {
		return domain.Rental{}, bad("%s hat bereits %d Artikel ausgeliehen (Maximum)", name, active)
	}
	if price <= 0 {
		price = s.prices[it.Type]
	}
	if start.IsZero() {
		start = time.Now()
	}
	r := &domain.Rental{ID: s.id(), ItemID: itemID, PersonID: personID, Start: start, Hours: hours, PricePerHour: price, SupervisorID: supervisorID}
	it.Rented = true
	s.rentals[r.ID] = r
	return *r, nil
}

func (s *Store) Return(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.rentals[id]
	if !ok || r.Returned {
		return bad("Ausleihe nicht gefunden oder bereits zurückgegeben")
	}
	r.Returned = true
	s.items[r.ItemID].Rented = false
	return nil
}

func (s *Store) Rentals() []map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rs := []domain.Rental{}
	for _, r := range s.rentals {
		rs = append(rs, *r)
	}
	sort.Slice(rs, func(a, b int) bool { return rs[a].ID > rs[b].ID })
	out := []map[string]any{}
	for _, r := range rs {
		out = append(out, r.View())
	}
	return out
}

func (s *Store) Prices() map[string]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := map[string]float64{}
	for k, v := range s.prices {
		out[k] = v
	}
	return out
}

func (s *Store) SetPrice(typ string, price float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.prices[typ]; !ok || price <= 0 {
		return bad("ungültiger Typ oder Preis")
	}
	s.prices[typ] = price
	return nil
}

func (s *Store) Stats() domain.Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st := domain.Stats{Total: len(s.items), Persons: len(s.persons)}
	for _, i := range s.items {
		if i.Rented {
			st.Rented++
		}
	}
	st.Available = st.Total - st.Rented
	return st
}
