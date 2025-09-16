package main

import (
	"github.com/bynow2code/geeksaver/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableTraverseRunHooks = true
	cmd.Execute()
}
