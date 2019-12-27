package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"io"
	"encoding/json"
	"sync"
	"sort"
	"strings"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
)

func TilesHandler(bucket *blob.Bucket) (http.HandlerFunc, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		path := req.URL.Path
		path = strings.TrimLeft(path, "/tiles")

		fh, err := bucket.NewReader(ctx, path, nil)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusNotFound)
			return
		}

		defer fh.Close()

		_, err = io.Copy(rsp, fh)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return 
		}

		fh.Close()
		return
	}

	return http.HandlerFunc(fn), nil
}

func CatalogHandler(bucket *blob.Bucket) (http.HandlerFunc, error) {

	mu := new(sync.RWMutex)
	images := make([]string, 0)
	
	var list_images func(context.Context, *blob.Bucket, string) error

	list_images = func(ctx context.Context, b *blob.Bucket, prefix string) error {
		
		iter := b.List(&blob.ListOptions{
			Delimiter: "/",
			Prefix:    prefix,
		})
		
		for {
			
			select {
			case <-ctx.Done():
				return nil
			default:
				// pass
			}
			
			obj, err := iter.Next(ctx)
			
			if err == io.EOF {
				break
			}
			
			if err != nil {
				return err
			}
			
			if obj.IsDir {
				
				err := list_images(ctx, b, obj.Key)
				
				if err != nil {
					return err
				}
				
				continue
			}
			
			fname := filepath.Base(obj.Key)
			
			if fname != "info.json" {
				return nil
			}
			
			id := filepath.Dir(obj.Key)				
			id = strings.TrimLeft(id, "/")
			
			mu.Lock()
			defer mu.Unlock()
			
			images = append(images, id)
			return nil
		}
		
		return nil
	}
	
	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()	
		
		err := list_images(ctx, bucket, "")

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

	ctx := context.Background()
	
	www_path, err := filepath.Abs(*www_root)

	if err != nil {
		log.Fatal(err)
	}

	tiles_bucket, err := blob.OpenBucket(ctx, *tiles_root)

	if err != nil {
		log.Fatal(err)
	}
	
	defer tiles_bucket.Close()
	
	tiles_handler, err := TilesHandler(tiles_bucket)

	if err != nil {
		log.Fatal(err)
	}
	
	catalog_handler, err := CatalogHandler(tiles_bucket)

	if err != nil {
		log.Fatal(err)
	}

	www_dir := http.Dir(www_path)
	www_handler := http.FileServer(www_dir)
	
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
