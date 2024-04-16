package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.Get("/", handler)

	r.Get("/about", aboutHandler)

	fmt.Println("Server is running on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	serveFile(w, r, "tmpl/index.tmpl")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	serveFile(w, r, "tmpl/about.tmpl")
}

func serveFile(w http.ResponseWriter, r *http.Request, filePath string) {
	file, err := os.Open(filePath)
	checkError(w, err)
	defer file.Close()

	fi, err := file.Stat()
	checkError(w, err)

	http.ServeContent(w, r, fi.Name(), fi.ModTime(), file)
}

func checkError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error opening file:", err)
		return
	}
}
