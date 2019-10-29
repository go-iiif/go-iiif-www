# go-iiif-www

A simple web application for testing go-iiif related tools.

## Important

This is work in progress. It is not properly documented yet. It is also, as of this writing, not a full-fledged application but rather a debugging tool for testing the output the of the `go-iiif` and `go-iiif-vips`. For example:

![](docs/images/tile-seed-bunk.png)

## Tools

### server

```
go run cmd/server/main.go -tiles-root /path/to/go-iiif-vips/docker/cache/ -www-root ./www/
2019/10/29 12:46:06 listening on localhost:8080
```

## See also

* https://github.com/go-iiif/go-iiif