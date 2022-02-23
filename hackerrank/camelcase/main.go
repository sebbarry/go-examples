package main

import (
	"fmt"
	"strings"
)

func main() {
    var input string //input string

    fmt.Scanf("%s\n", &input) //conver stdin file into a single string

    //b := []byte(input)

    /*
    when iterating through a string in go, 
    go will iterate through each and figure out what charactes are runes. 
    a rune is sometimes larger than a byte, so each character may actually increase its index by more than one.
    */
    var answer = 1
    for _, ch := range input {
        str := string(ch)
        if strings.ToUpper(str) == str {
            answer++
        }
        //min, max := 'A', 'Z'
        //if ch >= min && ch <= max {
            //it is a cpiatl letter.
         //   answer++
        //}
    }
    fmt.Println(answer)
}
