package main
import (
    "fmt"
)

func main() {

    fmt.Println("printing value..")
    v := 0
    for {
        if v <= 10 {
            go makeFile(&v)
            continue
        }
        break
    }
}


func makeFile(v *int){
    defer func () {
        fmt.Println("closing connection.")
        return
    }()
    temp := *v; //derefencing the value of thepointer v
    for i := 0; i < *v; i++ {
        fmt.Printf("Current Value: %i", temp)
    }
    temp++ // increase the currentvalue of temp
    *v = temp //put the value of tmep back into the v memory slot
}
