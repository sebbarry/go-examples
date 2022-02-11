package main

import (
    "flag"
    "os"
    "fmt"
    "encoding/csv"
    "strings"
    "time"
)

func usage() {
    fmt.Println("go build main.go && ./main -csv=<filename>")
    os.Exit(0)
}

func main() {

    //// define flags
    //csvfilename to use (problems.csv by default)
    csvFileName := flag.String("csv", "problems.csv", "a csv file in the format of question,answer")
    //time limit to use (30 seconds as default)
    timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")


    flag.Parse()

    file, err := os.Open(*csvFileName)
    if err != nil {
        usage()
        exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFileName))
    }
    r := csv.NewReader(file)     //make the new reader to read the file and output bytes
    lines, err := r.ReadAll()    //read all the lines
    if err != nil {
        exit("Failed to parse the provided CSV file.")
    }

    problems := parseLines(lines) //parse the lines from the problems.
    timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

    /*
        keep a counter of all correct values.
    */
    //start looping through each problem
    correct := 0
    problemloop: //this is a label - like an anchor we can reset the code to
    for i, p := range(problems) {
        fmt.Printf("Problem #%d: %s = \n", i+1, p.q) //read out the problem

        answerCh:=make(chan string) //make a channel to retrieve data from 
        //the go routine below.

        go func() {
            //we wnat to move this so taht is is non blocking
            //- put into a go routine
            var answer string
            fmt.Scanf("%s\n", &answer)     //read in the answer from th euser. 

            answerCh <- answer //when we get an answer, we send it into the 
            //answerCh channel we made above.
            //(the arrow points towards the way the data is moving.)

        }() //this is an iife


        //select statement: this is actually watching out for updates to 
        //channels we have defined.
        //we have a few cases we need to watch out for. 
        select {
        //if the time has ended...
        case <-timer.C: //if the timer.Channel has received a mess,
                        //timer has end.
            break problemloop
        //if we have received an answer from the answer channel
        case answer := <-answerCh:
            if answer == p.a {
                correct++
            }
        }
    }
    fmt.Printf("You scored: %d/%d", correct, len(problems))

}



func parseLines(lines [][]string) []problem {
    ret := make([]problem, len(lines))  //make a new []problem slice iwth length of lines
    for i, line := range lines {             //loop through each line in the csv file
        ret[i] = problem {                   //make a new problem
            q: line[0],                      //question is the first value in the array
            a: strings.TrimSpace(line[1]),   //answer is the second value in the array
        }
    }
    return ret;
}


type problem struct {
    q string
    a string
}



func exit(msg string) {
    fmt.Println(msg)
    os.Exit(1)
}
