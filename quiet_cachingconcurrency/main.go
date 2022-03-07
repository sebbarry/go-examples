package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gophercises/quiet_hn/hn"
)



//TYPES: 

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

type result struct {
    idx int //keep track of the index as we queue them up.
    item item
    err error
}

type storyCache struct {
    numStories int
    cache      []item
    expiration time.Time
    duration   time.Duration
    mutex      sync.Mutex
}

//global vars.
var (
    cache           []item
    cacheExpiration time.Time
    cacheMutex      sync.Mutex
)

//END Types

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
    //grab the number of stories we want ot acquire for the page.
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml")) //make a tepmlate to render

	http.HandleFunc("/", handler(numStories, tpl)) //pass the template and 
    //the number of storeis to our handler.

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}


/*
This handles our / route. 
*/
func handler(numStories int, tpl *template.Template) http.HandlerFunc {

    sc := storyCache{
        numStories: numStories,
        duration: 3 * time.Second,
    }

    //make sure to start this go routine whenever we call our handler.
    go func () {
        /*
        when this ticker goes off, go and populate whatever the cache should be.
        */
        ticker := time.NewTicker(3 * time.Second) //make a new ticker.
        for { //start while loop.
            temp := storyCache{
                numStories: numStories,
                duration: 6 * time.Second,
            }
            temp.stories() //handle querying the stories.
            //lock the mutex 
            sc.mutex.Lock()
            //update the cache
            sc.cache = temp.cache
            //reupdate the expiration
            sc.expiration = temp.expiration
            //unlock the mutex
            sc.mutex.Unlock()
            //
            <-ticker.C //output the channel time
        }
    }()


	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() //get initial time value of request
        stories, err := sc.stories()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
        err = tpl.Execute(w, data) //render the template
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}



func (sc *storyCache) stories() ([]item, error){ //make sure to always pass the story cache as a pointer to ensure unique
    sc.mutex.Lock()
    defer sc.mutex.Unlock() //remember to defer the mutex unlock.
    //check if the time is - meaning that our cache is stale
    if time.Now().Sub(sc.expiration) < 0 {
        return sc.cache, nil
    }
    stories, err := getTopStories(sc.numStories)
    if err != nil {
        return nil, err
    }
    sc.expiration = time.Now().Add(5 * time.Minute)
    sc.cache = stories
    return sc.cache, nil
}





func getTopStories(numStories int) ([]item, error) {
    var client hn.Client
    ids, err := client.TopItems()
    if err != nil {
        return nil, errors.New("Failed to return top stories.")
    }
    var stories []item
    at := 0
    //update the values so that we guarentee we get enough stories (minimum 30)
    for len(stories) < numStories {
        need := (numStories - len(stories)) * 5 / 4
        stories = append(stories, getStories(ids[at:at+need])...)
        at += need
    }
    return stories[:numStories], nil
}



func getStories(ids []int) []item {
    //had type result struct here.
    resultCh := make(chan result) //generate a channel with a result type


    //want := numStories * 5/4 //define the exact number of stories that we want.

    for i := 0; i < len(ids); i++ {
        go func(idx, id int) { //generate a go routine
            var client hn.Client //make a client for us to make syscalls with.
            hnItem, err := client.GetItem(id)
            if err != nil {
                resultCh <- result{err: err}
            }
            resultCh <- result{idx: idx, item: parseHNItem(hnItem)}
        }(i, ids[i])
    }

    //collect all the results[]
    var results []result
    //sort the results.
    for i:= 0; i < len(ids); i++ {
        results = append(results, <-resultCh)
    }
    sort.Slice(results, func(i, j int) bool { //sort the results by idx value.
        return results[i].idx < results[j].idx
    })

    var stories []item
    for _, res := range results {
        if res.err != nil {
            continue
        }
        if isStoryLink(res.item) {
            stories = append(stories, res.item)
        }
    }
    return stories
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}
