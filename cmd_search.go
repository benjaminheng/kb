package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func NewSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search",
		Short: "Searches contents of the knowledge base",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			search(args)
		},
	}
}

func search(args []string) {
	// search contents of all files in kb
	fmt.Println("Not implemented yet")
	r := strings.NewReader("hello")
	b := &bytes.Buffer{}
	runShellCommand("fzf", r, b)
	fmt.Printf("b = %+v\n", b)
}
