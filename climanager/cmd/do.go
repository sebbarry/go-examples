package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)


var doCmd = &cobra.Command{
    Use: "do",
    Short: "do commands.",
    Long: "long" ,
    Run: func(cmd *cobra.Command, args []string) {
        var ids []int
        for _, arg := range args {

            //Atoi conversion to convert string values to int.
            //more precisely, ASCII string to integer
            id, err := strconv.Atoi(arg)

            if err != nil {
                fmt.Println("failed to parse arg:", arg)
            } else {
                ids = append(ids, id)
            }
        }
        fmt.Println(ids)
    },
}



func init() {
    RootCmd.AddCommand(doCmd)
}
