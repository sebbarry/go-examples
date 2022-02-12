package main


import (
    "fmt"
    "net/http"
    "urlshort"
)


func main() {

    mux := defaultMux()                     //make a default mux

    pathsToUrls := map[string]string {      //make a map of urls.
        "/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
        "/yaml-godoc":     "https://godoc.org/goplkg.in/yaml.v2",
    }

    mapHandler := urlshort.MapHandler(pathsToUrls, mux) //this is the map handler

    yaml := `
    - path: /urlshort
      url: https://github.com/gophercises/urlshort
    - path: /urlshort-final
      url: https://github.com/gophercises/urlshort/tree/final
    `
    yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
    if err != nil {
        panic(err)
    }

    /*
    we can pass a deafult route to use (below) in case we do not find a route to use.
    */

    fmt.Println("Starting to serve on :3000")
    http.ListenAndServe(":3000", yamlHandler)


}


//this is the default router to use.
func defaultMux() *http.ServeMux {
    mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}


//this is a defualt response.
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

