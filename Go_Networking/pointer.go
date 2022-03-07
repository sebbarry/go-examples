package main;


import (
    "fmt"
)

type Node struct {
    val int;
    next Node
}

func main() {
    v := 5
    p := &v
    fmt.Print("this is the memory address of the pointer: ", p, "\n") //memory address of the value
    fmt.Print("this is the pointer value: ", *p)

    t := Node{5, nil}
    fmt.Println(t.val)
}











