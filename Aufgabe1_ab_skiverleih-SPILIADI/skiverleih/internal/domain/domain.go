package domain

import (
	"errors"
	"math"
	"time"
)

var ErrInvalid = errors.New("invalid input")

// Business rule limits (easy to change).
const (
	MaxHours          = 72 // longest allowed rental
	MaxItemsPerPerson = 4  // active items per person at the same time
	AdultAge          = 18 // below this age a registered adult supervisor is required
)

type Item struct {
	ID     int    `json:"id"`
	Number string `json:"number"`
	Type   string `json:"type"`
	Rented bool   `json:"rented"`
}

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Rental struct {
	ID           int       `json:"id"`
	ItemID       int       `json:"item_id"`
	PersonID     int       `json:"person_id"`
	Start        time.Time `json:"start"`
	Hours        int       `json:"hours"`
	PricePerHour float64   `json:"price_per_hour"`
	Returned     bool      `json:"returned"`
	SupervisorID int       `json:"supervisor_id"` // required if the person is under AdultAge
}

func (r Rental) End() time.Time { return r.Start.Add(time.Duration(r.Hours) * time.Hour) }
func (r Rental) Total() float64 {
	return math.Round(float64(r.Hours)*r.PricePerHour*(1-Discount(r.Hours))*100) / 100
}

// MarshalJSON adds the computed fields end/total.
func (r Rental) View() map[string]any {
	return map[string]any{"id": r.ID, "item_id": r.ItemID, "person_id": r.PersonID, "start": r.Start,
		"end": r.End(), "hours": r.Hours, "price_per_hour": r.PricePerHour, "total": r.Total(), "returned": r.Returned, "supervisor_id": r.SupervisorID,
		"discount": Discount(r.Hours), "subtotal": float64(r.Hours) * r.PricePerHour}
}

type Stats struct {
	Available int `json:"available"`
	Rented    int `json:"rented"`
	Total     int `json:"total"`
	Persons   int `json:"persons"`
}

// Tier is a duration discount: from MinHours on, Discount (0.1 = 10 %) is deducted from the total.
type Tier struct {
	MinHours int     `json:"min_hours"`
	Discount float64 `json:"discount"`
}

var Tiers = []Tier{{4, 0.10}, {8, 0.20}}

// Discount returns the highest discount that applies to a rental of the given hours.
func Discount(hours int) float64 {
	d := 0.0
	for _, t := range Tiers {
		if hours >= t.MinHours && t.Discount > d {
			d = t.Discount
		}
	}
	return d
}
