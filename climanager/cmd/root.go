package cmd

import (
    "github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
    Use: "task",
    Short: "Hugo is a very fast static site genrrator",
}

