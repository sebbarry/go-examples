package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", recoverMw(mux, true)))
}

//this is a middleware function that we pass as the server's first starting
//function
func recoverMw(app http.Handler, dev bool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            err := recover()
            if err != nil {
                log.Println(err)
                stack := debug.Stack()
                log.Println(string(stack))
                if dev == false {
                    http.Error(w, "Something wnet wrong: ", http.StatusInternalServerError);
                    return
                }
                w.WriteHeader(http.StatusInternalServerError) //return 500 sc in header.
                fmt.Fprintf(w, "<h1>Pannic: %v</h1><pre>%s</pre>", err, string(stack))
            }
        }()

        nw := &responseWriter{ResponseWriter: w}
        app.ServeHTTP(w, r) //if anything in this function panics, the deferred function get called.
        nw.flush()
    }
}

//defering functions when panic'd are good for logging etc.
func panicDemo(w http.ResponseWriter, r *http.Request) {
    //when we panic, we always want a defer function
    //this function will get executed whenever the function exits.
    //this is called before the panicDemo function finishes.
    /*
    defer func() {
        //!! NOTE.
        err := recover() //recover function is a built in function...will not
        //do anything useful unless we are in a defer function.
        fmt.Fprint(w, err)
    }()
    */
	funcThatPanics()
}


type responseWriter struct {
    http.ResponseWriter
    writes [][]byte
    status int
}

func (rw *responseWriter) Write(b []byte) (int, error)  {
    rw.writes = append(rw.writes, b)
    return len(b), nil
}


func (rw *responseWriter) flush() error {
    if rw.status != 0 {
        rw.ResponseWriter.WriteHeader(rw.status)
    }
    for _, write := range rw.writes {
        _, err := rw.ResponseWriter.Write(write)
        if err != nil {
            return err
        }
    }
    return nil
}


func (rw *responseWriter) WriteHeader(statusCode int) {
    rw.status = statusCode
}


func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
    hijacker, ok := rw.ResponseWriter.(http.Hijacker) //wrapping an interface
    if !ok {
        return nil, nil, fmt.Errorf("The responsewriter deos not support the hijacker interface")
    }
    return hijacker.Hijack()
}

func (rw *responseWriter) Flush() {
    flusher, ok := rw.ResponseWriter.(http.Flusher)
    if !ok {
        return
    }
    flusher.Flush()
    return
}


func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("panic")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
