package main

import (
	"log"
	"io/ioutil"
	"net/http"
	"github.com/shurcooL/vfsgen"
)

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
	var fs http.FileSystem = http.Dir("./templates")

	err := vfsgen.Generate(fs, vfsgen.Options{})
	if err != nil {
		log.Fatalln(err)
	}
	
	http.HandleFunc("/", indexHandler)
	log.Println("starting app...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
