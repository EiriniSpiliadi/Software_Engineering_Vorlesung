package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"skiverleih/internal/domain"
	"skiverleih/internal/store"
)

type Handler struct{ S *store.Store }

func write(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func fail(w http.ResponseWriter, err error) { http.Error(w, err.Error(), http.StatusBadRequest) }

func pathID(r *http.Request) int { n, _ := strconv.Atoi(r.PathValue("id")); return n }

func (h *Handler) Routes() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("GET /api/dashboard", func(w http.ResponseWriter, r *http.Request) { write(w, h.S.Stats()) })
	m.HandleFunc("GET /api/items", func(w http.ResponseWriter, r *http.Request) { write(w, h.S.Items()) })
	m.HandleFunc("POST /api/items", func(w http.ResponseWriter, r *http.Request) {
		var in struct{ Number, Type string }
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			fail(w, err)
			return
		}
		it, err := h.S.AddItem(in.Number, in.Type)
		if err != nil {
			fail(w, err)
			return
		}
		write(w, it)
	})
	m.HandleFunc("DELETE /api/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		if err := h.S.DeleteItem(pathID(r)); err != nil {
			fail(w, err)
			return
		}
		write(w, true)
	})
	m.HandleFunc("GET /api/persons", func(w http.ResponseWriter, r *http.Request) { write(w, h.S.Persons()) })
	m.HandleFunc("POST /api/persons", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Name string
			Age  int
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			fail(w, err)
			return
		}
		p, err := h.S.AddPerson(in.Name, in.Age)
		if err != nil {
			fail(w, err)
			return
		}
		write(w, p)
	})
	m.HandleFunc("GET /api/rentals", func(w http.ResponseWriter, r *http.Request) { write(w, h.S.Rentals()) })
	m.HandleFunc("POST /api/rentals", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			ItemID       int       `json:"item_id"`
			PersonID     int       `json:"person_id"`
			Start        time.Time `json:"start"`
			Hours        int       `json:"hours"`
			PricePerHour float64   `json:"price_per_hour"`
			SupervisorID int       `json:"supervisor_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			fail(w, err)
			return
		}
		rt, err := h.S.Rent(in.ItemID, in.PersonID, in.Start, in.Hours, in.PricePerHour, in.SupervisorID)
		if err != nil {
			fail(w, err)
			return
		}
		write(w, rt.View())
	})
	m.HandleFunc("POST /api/rentals/{id}/return", func(w http.ResponseWriter, r *http.Request) {
		if err := h.S.Return(pathID(r)); err != nil {
			fail(w, err)
			return
		}
		write(w, true)
	})
	m.HandleFunc("GET /api/tiers", func(w http.ResponseWriter, r *http.Request) { write(w, domain.Tiers) })
	m.HandleFunc("GET /api/prices", func(w http.ResponseWriter, r *http.Request) { write(w, h.S.Prices()) })
	m.HandleFunc("PUT /api/prices/{type}", func(w http.ResponseWriter, r *http.Request) {
		var in struct{ Price float64 }
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			fail(w, err)
			return
		}
		if err := h.S.SetPrice(r.PathValue("type"), in.Price); err != nil {
			fail(w, err)
			return
		}
		write(w, true)
	})
	return m
}
