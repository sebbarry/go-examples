package main 

import (
    "fmt"
)

func main() {
    var length, delta int
    var input string
    fmt.Scanf("%d\n", &length)
    fmt.Scanf("%s\n", &input)
    fmt.Scanf("%d\n", &delta)


    var ret []rune
    for _, ch := range input {
        ret = append(ret, cipher(ch, delta))
    }
    fmt.Println(string(ret))

}

func cipher(r rune, delta int) rune {
    if r >= 'A' && r <= 'Z' {
        return rotateWithBase(r, 'A', delta)
    }
    if r >= 'a' && r <= 'z' {
        return rotateWithBase(r, 'a', delta)
    }
    return r
}

func rotateWithBase(r rune, base, delta int) rune {
    tmp := int(r) - base
    tmp = (tmp + base) % 26
    return rune(tmp + base)
}



//rotate function rotates the strings by the anticipated rotation 
//delta

/*
func rotate(s rune, delta int, key []rune) rune {
    idx := -1
    for i, r := range key {
        if r == s {
            idx = i
            break
        }
    }
    if idx < 0 {
        panic("idx < 0")
    }
    //elegent way to find the rotation key
    idx = (idx + delta) % len(key)
    return key[idx]
}
*/
