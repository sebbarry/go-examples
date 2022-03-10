/*
README
this file is the server.
We send requests to this server - it packs the data into a cache over an n-second time period - finally flushing the data over a network request to another server.
*/


package main

import (
    "fmt"
    "net"
    "time"
    "sync"
    "flag"
)

const (
    CONN_PORT = "8080"
    CONN_SERVER = "localhost"
    CONN_TYPE = "tcp"
)

var (
    DEST_SERVER string
    DEST_PORT string
)


const (
   CACHE_DURATION = 8000 //mS
)


type Cache struct {
    data []byte
}

type Timer struct {
    time time.Time
}



func (c *Cache) flush() {
    //handle flush here.
    //fmt.Println(string(c.data))
    err := c.transferChunks()
    if err != nil {
        fmt.Println(err)
    }
    c.data = nil
}


func (c *Cache) transferChunks() error {
    //make network connection here and flush the queue here 
    fmt.Println("transfering chunks now.")
    if cap(c.data) == 0 {
        return nil
    }
    l := (cap(c.data)-1)/10 //upper bound / 10 slices
    for i := 0; i < 9; i++ {
        go func(beg int, end int, c []byte) {
            d := c[beg:end] //grab a chunk of the byte array
            conn, err := net.Dial("tcp", DEST_SERVER + ":" + DEST_PORT)
            if conn == nil {
                return
            } else{
                defer conn.Close()
            }
            if err != nil {
                fmt.Println(err)
                return
            }
            _, err = conn.Write(d)
            if err != nil {
                fmt.Println(err)
            }
        }(i*l, (i+1)*l, c.data)
    }
    return nil
}


var wg sync.WaitGroup


func main() {
    //get flags
    ip := *flag.CommandLine.String("ip", "localhost", "destination server")
    port := *flag.CommandLine.String("port", "8000", "destination port")
    flag.Parse()
    if len(ip) == 0 || len(port) == 0 {
        panic("Please include a destination server ip and port")
    }
    DEST_SERVER = ip //assign destination server from the cli flag.
    DEST_PORT = port

    l, err := net.Listen(CONN_TYPE, CONN_SERVER + ":" + CONN_PORT)
    if err != nil {
        panic(err)
    }

    //make new vars.
    cache := Cache{make([]byte, 1024)}
    timer := Timer{time.Now()}
    //end make
    go listen(l, &cache) //start our loop for the server.
    go startTimer(&timer, &cache)
    wg.Add(1)
    wg.Wait()
}


func startTimer(t *Timer, c *Cache) {
    //handle the timer elapsation here.
    wg.Add(1)
    for {
        now := time.Now()
        if now.Sub(t.time).Milliseconds() > CACHE_DURATION {
            //handle duration elapse
            c.flush()
            t.time = now
        }
    }
}


//handle listening of the server.
func listen(l net.Listener, c *Cache) {
    wg.Add(1)
    fmt.Println("Listening to: " + CONN_SERVER + ":" + CONN_PORT)
    for {
        conn, err := l.Accept()
        if err != nil {
            panic(err)
        }
        //handle the connection.
        go handleCon(conn, *(&c))
    }
}


func handleCon(conn net.Conn, cache *Cache) {
    defer func() {
        fmt.Println("Closing connection: ", conn.LocalAddr().String())
        conn.Close()
    }()
    for {
        buf := make([]byte, 1024) //make buffer for data input
        data, err := conn.Read(buf) //read data into buffer.
        if err != nil {
            break
        }
        if data <= 1 {
            conn.Close()
            break
        }
        //append the data to the message cache.
        cache.data = append(cache.data, buf...) //spread this guy
    }
}


