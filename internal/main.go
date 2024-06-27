package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jayagr26/url-shortner/internal/urlshort"
)

func main() {
	yamlFile := flag.String("yml", "url.yml", "path to yaml file")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

// 	// Build the YAMLHandler using the mapHandler as the
// 	// fallback
	ymlData, err := os.ReadFile(*yamlFile) 
	if err != nil {
		log.Fatalln("Unable to read yml data: ", err)
	}
	yamlHandler, err := urlshort.YAMLHandler(ymlData, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}