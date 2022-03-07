package main

import (
    "fmt"
    secret "cli_secrets"
)

func main() {
    v := secret.File("fake-key", ".secrets")
    err := v.Set("demo_key1", "123 some value")
    if err != nil {
        panic(err)
    }
    plain, err := v.Get("demo_key1")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Value: ", plain)
}
