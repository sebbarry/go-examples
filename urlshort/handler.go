package urlshort


import (
    "net/http"
    "gopkg.in/yaml.v2"
)


//takes in values and searches to see if we have a matching key/value
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

    return func(w http.ResponseWriter, r *http.Request) {
        path := r.URL.Path               //take the path from the url request.
        if dest, ok := pathsToUrls[path]; ok { //ifwe find a key in teh map, 
            http.Redirect(w, r, dest, http.StatusFound)
            return
        }
        fallback.ServeHTTP(w, r)
    }
}



func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
    //1. parsethe yamlsomehow
    pathUrls, err := parseYaml(yamlBytes)

    if err != nil {
        return nil, err
    }

    //2. conver yaml array into map type
    pathsToUrls := buildMap(pathUrls)

    //3. return the maphandler endpoint
    return MapHandler(pathsToUrls, fallback), nil
}


//parse the yaml file
func parseYaml(data []byte) ([]pathURL, error) {
    var pathURLs []pathURL                          //define an array of pathURLs
    err := yaml.Unmarshal(data, &pathURLs)          //this unmarshals the yaml bytes.
    if err != nil {
        return nil, err
    }
    return pathURLs, nil
}


//build a map structure to return to the user.
func buildMap(pathUrls []pathURL) map[string]string {
    pathsToUrls := make(map[string]string)              //make a map: "string":"string"
    for _, pu := range pathUrls {                       //loop trough each value in the map
        pathsToUrls[pu.Path] = pu.URL
    }
    return pathsToUrls
}

//struct to hold path and url types.
type pathURL struct {
    //tags to the right are automatically parsed from data that is being passed in from the .yaml definition.
    Path string `yaml:"path"`   //tag: yaml:path
    URL string `yaml:"url"`     //tag: yaml:url
}
