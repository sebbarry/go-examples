package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

/*
add args here.
*/
var addCmd = &cobra.Command{
    Use: "add",
    Short: "Adds a task to your task list.",
    Run: func(cmd *cobra.Command, args []string) {
        task := strings.Join(args, " ")
        fmt.Printf("Added \"%s\" to your task list.\n", task)
    },
}

/*
init is a function taht can be run before main. the init functions are alwaays priority functins that can set things up for us before the actual main() function of a go program is run.
*/
func init() {
    RootCmd.AddCommand(addCmd)
}
