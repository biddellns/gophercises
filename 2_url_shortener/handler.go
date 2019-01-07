package urlshort

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

		if dest, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	})
}


type UrlData struct {
	Path string `yaml:"path"`
	Url string `yaml:"url"`
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	urlData, err := parseYaml(yml)

	if err != nil {
		return nil, err
	}

	mappedUrlData := getMappedUrlData(urlData)

	return MapHandler(mappedUrlData, fallback), nil
}

func parseYaml(yml []byte) ([]UrlData, error) {
	var urlData []UrlData

	err := yaml.Unmarshal(yml, &urlData)

	return urlData, err
}

func getMappedUrlData(urlData []UrlData) (map[string]string) {
	mappedUrlData := make(map[string]string, len(urlData))

	for _, item := range urlData {
		mappedUrlData[item.Path] = item.Url
	}

	return mappedUrlData
}