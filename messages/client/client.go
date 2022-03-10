package main

import (
    "fmt"
    "net"
)

func main() {
    l, err := net.Listen("tcp", "localhost" + ":" + "8000")
    if err != nil {
        panic(err)
    }
    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println(err)
        }
        go handleCon(conn)
    }
}

func handleCon(n net.Conn) {
    fbuf := make([]byte, 1024)
    for {
        tbuf := make([]byte, 1024)
        d, err := n.Read(tbuf)
        if err != nil {
            break
        }
        if d <= 1 {
            n.Close()
            break
        }
        fbuf = append(fbuf, tbuf...)
    }
    fmt.Print(string(fbuf))
}

