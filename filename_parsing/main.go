/*
This is just a writeup script. Throwaway code...
*/

package main

import (
	"fmt"
	"os"
)

func main() {
    var filename string = "birthday_001.txt"
    // => Birthday - 1 of 4.txt <- rename to this
    newName, err := match(filename)
    if err != nil {
        fmt.Println("no match")
        os.Exit(1)
    }
    fmt.Println(newName)
}


// match returns the new filename or an error if the file name 
//didnt match the pattern.
func match(filename string) (string, error) {
    return "", nil;
}
