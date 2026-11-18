package main

import (
	"fmt"
	"log"
	"net/http"

	"fahrradverleih/internal/handler"
	"fahrradverleih/internal/store"
)

func main() {
	types := []string{ "City", "Mountain", "E-Bike" }
	s := store.New(types, 5)
	s.SetPrice("City", 3)
	s.SetPrice("Mountain", 4)
	s.SetPrice("E-Bike", 6)
	for i := 1; i <= 6; i++ { // Demo-Daten
		s.AddItem(fmt.Sprintf("BIKE-%03d", i), types[i%len(types)])
	}
	s.AddPerson("Anna Meister", 25)
	s.AddPerson("Ben Stiller", 16)
	s.AddPerson("Eirini Spiliadi", 20)

	mux := (&handler.Handler{S: s}).Routes()
	mux.Handle("/", http.FileServer(http.Dir("web")))
	log.Println("Fahrrad-Verleih läuft auf http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
