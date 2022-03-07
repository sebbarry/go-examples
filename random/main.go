package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)


var (
    mutex sync.Mutex

    /*
    This is to allow us to wait for multiple go routines to finish their concurrent operations before we move on to the next thing.
    We use waitgroups instead of channels, if we just want the go routines to run w/o knowing the results of each function call. We could use a channel, however that would only be useful if we wanted to retrieve some result of the go routine.
    */
    wg sync.WaitGroup
    links []string
    data  []string
)

func main() {
    populateLinks(&links) //pass in the location of the links list.
    for _, urls := range links {
        wg.Add(1) //add a new value to the waitgroup
        go func(urls string) {
            defer wg.Done()
            res, err := http.Get(urls)
            fmt.Println(res.Status)
            if err != nil {
                fmt.Println(err)
                return
            }
            if res.StatusCode < 200 || res.StatusCode > 299 {
                fmt.Printf("invalid status code: %i", res.StatusCode)
                return
            }
            buf := new(bytes.Buffer)
            buf.ReadFrom(res.Body)
            body := buf.String()
            fmt.Println("locking...")
            mutex.Lock()
            data = append(data, body)
            mutex.Unlock()
            fmt.Println("unlocking...")
        }(urls)
    }
    wg.Wait()
    fmt.Println("done")
}



func populateLinks(list *[]string) {
    urls := []string{
        "http://google.com",
        "http://youtube.com",
        "http://facebook.com",
        "http://cbaminvestments.com",
        "http://sms-inc.net",
        //populate other links here, 
    }
    for _, i := range urls {
        *list = append(*list, i)
    }
}
