package main


import (
    "fmt"
    "strings"
    "os"

    "link_parser"
)





/*
var exampleHtml = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">A link to another page</a>
</body>
</html>
`
*/


func main() {
    //loop through all the .html files in the files directory
    files, err := os.ReadDir("files/")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }


    //make a channel
    ch := make(chan int)

    for _, file := range files {
        f, err := os.Open("files/"+file.Name()) //open the file
        if err != nil {
            fmt.Printf("error opening file: %s. \n", file.Name())
            fmt.Printf("error: %s\n", err)
            fmt.Println("-----")
            continue
        }
        defer f.Close() //close the file when done

        fsize, err := f.Stat()     //get stat of file.
        if err != nil {
            fmt.Printf("error getting file stat: %s\n", file.Name())
            continue
        }

        size := fsize.Size()       //get the file's size to load into memeory buffer.
        data := make([]byte, size) //make an empty buffer
        count, err := f.Read(data) //read file into buffer.

        if err != nil {
            fmt.Printf("error reading data into buffer\n")
            continue
        }
        //fmt.Printf("read %d bytes: %q\n", count, data[:count])

        go parseData(string(data[:count]), &ch)

    }

    //we want to make sure the channel closes after our go routines have finished running.
    d:=0
    for {
        select{
            case dd := <-ch:
            {
                d += dd
                if d >= len(files) {
                    close(ch)
                    os.Exit(0)
                }
            }
        }
    }

}


func parseData(s string, c *chan int) {
    r := strings.NewReader(s) //make a new reader to read a string.
    links, err := link_parser.Parse(r)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", links)
    *c <- 1
}
