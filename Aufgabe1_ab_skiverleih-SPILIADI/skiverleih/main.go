package main

import (
	"fmt"
	"log"
	"net/http"

	"skiverleih/internal/handler"
	"skiverleih/internal/store"
)

func main() {
	types := []string{ "Alpin", "Langlauf", "Snowboard" }
	s := store.New(types, 5)
	s.SetPrice("Alpin", 8)
	s.SetPrice("Langlauf", 5)
	s.SetPrice("Snowboard", 7)
	for i := 1; i <= 6; i++ { // Demo-Daten
		s.AddItem(fmt.Sprintf("SKI-%03d", i), types[i%len(types)])
	}
	s.AddPerson("Rita Schillt", 15)
	s.AddPerson("Dimitris Papadopoulos", 30)
	s.AddPerson("Eirini Spiliadi", 20)

	mux := (&handler.Handler{S: s}).Routes()
	mux.Handle("/", http.FileServer(http.Dir("web")))
	log.Println("Ski-Verleih läuft auf http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
