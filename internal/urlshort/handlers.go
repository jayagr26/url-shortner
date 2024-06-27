package urlshort

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(ymlData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var paths []map[string]string

	err := yaml.Unmarshal(ymlData, &paths)
	if err != nil {
		log.Println("Error occurred while parsing yaml data: ", ymlData)
		return nil, err
	}
	
	pathsToUrls := make(map[string]string)
	
	for _, path := range paths {
		pathsToUrls[path["path"]] = path["url"]
	}

	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var paths []map[string]string

	if err := json.Unmarshal(jsonData, &paths); err != nil {
		log.Println("Error occurred while parsing json data: ", jsonData)
		return nil, err
	}

	pathsToUrls := make(map[string]string)

	for _, path := range paths {
		pathsToUrls[path["path"]] = path["url"]
	}

	return MapHandler(pathsToUrls, fallback), nil
}