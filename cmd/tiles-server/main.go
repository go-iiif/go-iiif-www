package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"os"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"encoding/json"
	"sync"
	"sort"
)

func CatalogHandler(tile_root string) (http.HandlerFunc, error) {

	mu := new(sync.RWMutex)
	
	fn := func(rsp http.ResponseWriter, req *http.Request) {

		images := make([]string, 0)
		
		cb := func(path string, info os.FileInfo) error {
		
			if info.IsDir() {
				return nil
			}
			
			fname := filepath.Base(path)

			if fname != "info.json" {
				return nil
			}

			root := filepath.Dir(path)
			image := filepath.Base(root)

			mu.Lock()
			defer mu.Unlock()

			images = append(images, image)
			return nil
		}
		
		cr := crawl.NewCrawler(tile_root)
		err := cr.Crawl(cb)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return 
		}

		sort.Strings(images)

		enc, err := json.Marshal(images)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return 
		}

		rsp.Write(enc)
	}

	return http.HandlerFunc(fn), nil	
}

func main() {

	host := flag.String("host", "localhost", "Hostname to listen on")
	port := flag.Int("port", 8080, "Port to listen on")

	tiles_root := flag.String("tiles-root", "", "")
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

	catalog_handler, err := CatalogHandler(tiles_path)

	if err != nil {
		log.Fatal(err)
	}
	
	mux := http.NewServeMux()

	mux.Handle("/catalog/", catalog_handler)	
	mux.Handle("/tiles/", tiles_handler)
	mux.Handle("/", www_handler)

	address := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("listening on %s\n", address)

	err = http.ListenAndServe(address, mux)

	if err != nil {
		log.Fatal(err)
	}

}
