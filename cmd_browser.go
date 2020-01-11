package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewBrowseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "browser",
		Short: "Opens your knowledge base in a web browser",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Not implemented yet")
		},
	}
}
