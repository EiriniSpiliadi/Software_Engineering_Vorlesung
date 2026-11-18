package store

import (
	"testing"
	"time"

	"fahrradverleih/internal/domain"
)

func TestRentFlow(t *testing.T) {
	s := New([]string{"A"}, 5)
	it, _ := s.AddItem("X-1", "A")
	p, _ := s.AddPerson("Anna", 30)
	r, err := s.Rent(it.ID, p.ID, time.Now(), 3, 0, 0)
	if err != nil || r.Total() != 15 {
		t.Fatalf("rent failed: %v total=%v", err, r.Total())
	}
	if _, err := s.Rent(it.ID, p.ID, time.Now(), 1, 0, 0); err == nil {
		t.Fatal("double rent must fail")
	}
	if s.Stats().Rented != 1 {
		t.Fatal("expected 1 rented")
	}
	if err := s.Return(r.ID); err != nil || s.Stats().Available != 1 {
		t.Fatal("return failed")
	}
}

func TestOneTypeAtATime(t *testing.T) {
	s := New([]string{"A", "B"}, 5)
	a1, _ := s.AddItem("A-1", "A")
	a2, _ := s.AddItem("A-2", "A")
	b1, _ := s.AddItem("B-1", "B")
	p, _ := s.AddPerson("Anna", 30)
	r1, _ := s.Rent(a1.ID, p.ID, time.Now(), 1, 0, 0)
	if _, err := s.Rent(b1.ID, p.ID, time.Now(), 1, 0, 0); err == nil {
		t.Fatal("other type must be rejected while a rental is active")
	}
	r2, err := s.Rent(a2.ID, p.ID, time.Now(), 1, 0, 0)
	if err != nil {
		t.Fatal("same type must be allowed")
	}
	s.Return(r1.ID)
	s.Return(r2.ID)
	if _, err := s.Rent(b1.ID, p.ID, time.Now(), 1, 0, 0); err != nil {
		t.Fatal("other type must be allowed after returning everything")
	}
}

func TestOverdueBlocksNewRental(t *testing.T) {
	s := New([]string{"A"}, 5)
	a1, _ := s.AddItem("A-1", "A")
	a2, _ := s.AddItem("A-2", "A")
	p, _ := s.AddPerson("Anna", 30)
	s.Rent(a1.ID, p.ID, time.Now().Add(-3*time.Hour), 1, 0, 0)
	if _, err := s.Rent(a2.ID, p.ID, time.Now(), 1, 0, 0); err == nil {
		t.Fatal("overdue rental must block new rentals")
	}
}

func TestLimits(t *testing.T) {
	s := New([]string{"A"}, 5)
	p, _ := s.AddPerson("Anna", 30)
	if _, err := s.AddPerson("anna", 30); err == nil {
		t.Fatal("duplicate person name must fail")
	}
	for i := 0; i < 5; i++ {
		it, _ := s.AddItem(string(rune('a'+i)), "A")
		_, err := s.Rent(it.ID, p.ID, time.Now(), 1, 0, 0)
		if (i < 4) != (err == nil) {
			t.Fatalf("item %d: unexpected result %v", i, err)
		}
	}
	x, _ := s.AddItem("z", "A")
	q, _ := s.AddPerson("Ben", 30)
	if _, err := s.Rent(x.ID, q.ID, time.Now(), 73, 0, 0); err == nil {
		t.Fatal("more than 72 hours must fail")
	}
}

func TestMinorNeedsSupervisor(t *testing.T) {
	s := New([]string{"A"}, 5)
	it, _ := s.AddItem("A-1", "A")
	kid, _ := s.AddPerson("Kid", 15)
	teen, _ := s.AddPerson("Teen", 17)
	adult, _ := s.AddPerson("Eva", 40)
	if _, err := s.Rent(it.ID, kid.ID, time.Now(), 1, 0, 0); err == nil {
		t.Fatal("minor without supervisor must fail")
	}
	if _, err := s.Rent(it.ID, kid.ID, time.Now(), 1, 0, teen.ID); err == nil {
		t.Fatal("minor as supervisor must fail")
	}
	if _, err := s.Rent(it.ID, kid.ID, time.Now(), 1, 0, adult.ID); err != nil {
		t.Fatal(err)
	}
}

func TestPricing(t *testing.T) {
	cases := []struct {
		hours int
		want  float64
	}{{1, 5}, {2, 10}, {3, 15}, {4, 18}, {8, 32}, {24, 96}} // 5 EUR/h, 10 % from 4 h, 20 % from 8 h
	for _, c := range cases {
		r := domain.Rental{Hours: c.hours, PricePerHour: 5}
		if r.Total() != c.want {
			t.Fatalf("%d h: got %v, want %v", c.hours, r.Total(), c.want)
		}
	}
}
