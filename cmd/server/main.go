package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

func main() {

	host := flag.String("host", "localhost", "Hostname to listen on")
	port := flag.Int("port", 8080, "Port to listen on")

	tiles_root := flag.String("tiles-root", "", "")
	// tiles_url := flag.String("tiles-url", "/tiles", "")

	www_root := flag.String("www-root", "", "")

	flag.Parse()

	tiles_path, err := filepath.Abs(*tiles_root)

	if err != nil {
		log.Fatal(err)
	}

	www_path, err := filepath.Abs(*www_root)

	if err != nil {
		log.Fatal(err)
	}

	tiles_dir := http.Dir(tiles_path)
	tiles_handler := http.FileServer(tiles_dir)
	tiles_handler = http.StripPrefix("/tiles/", tiles_handler)
	
	www_dir := http.Dir(www_path)
	www_handler := http.FileServer(www_dir)

	mux := http.NewServeMux()

	mux.Handle("/tiles/", tiles_handler)	
	mux.Handle("/", www_handler)

	address := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("listening on %s\n", address)

	err = http.ListenAndServe(address, mux)

	if err != nil {
		log.Fatal(err)
	}

}
