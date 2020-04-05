package main

import (
	"context"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

//go:generate go run -mod=vendor vfs.go

func indexHandler(w http.ResponseWriter, r *http.Request) {
	vars := make(map[string]string)
	vars["file"] = strings.TrimPrefix(r.URL.Path, "/")

	log.Println(vars["file"])

	file, err := assets.Open(vars["file"])
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`Error getting ` + vars["file"]))
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`Could not read ` + vars["file"]))
	}

	t := template.New("html")
	t, err = t.Parse(string(data))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`Could not read ` + vars["file"]))
	}
	t.Execute(w, t)
}

func main() {
	mux := http.NewServeMux()
	server := http.Server{Addr: ":8080", Handler: mux}

	var useTLS bool
	var tlsServer *http.Server
	tlsCert, tlsKey := os.Getenv("TLS_CERT"), os.Getenv("TLS_KEY")
	if tlsCert != "" && tlsKey != "" {
		useTLS = true
		tlsServer = &http.Server{Addr: ":8443", Handler: mux}
	} else {
		useTLS = false
		tlsServer = nil
	}

	mux.HandleFunc("/", indexHandler)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if useTLS {
		go func() {
			log.Print("https server started on port 8443")
			if err := tlsServer.ListenAndServeTLS(tlsCert, tlsKey); err != nil && err != http.ErrServerClosed {
				log.Fatalf("https listen failed: %s\n", err)
			}
		}()
	}

	go func() {
		log.Print("http server started on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http listen failed: %s\n", err)
		}
	}()

	<-done

	// Attempt to gracefully stop the server on interupt
	log.Print("server stopping...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	// Attempt to gracefully shutdown the server with a timeout
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Print("server shutdown complete")
}
