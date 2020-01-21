package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

//go:generate go run -mod=vendor vfs.go

func indexHandler(w http.ResponseWriter, r *http.Request) {
	file, err := assets.Open("/index.html")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`Error getting index.html`))
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`Could not read index.html`))
	}
	w.Write(data)
}

func main() {
	http.HandleFunc("/", indexHandler)
	log.Println("starting app...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
