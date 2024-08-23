/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var version = [4]int{0, 1, 0}
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "App version",
	Long:  `App version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(
			"tSFInstaller, version",
			strings.Trim(
				strings.Join(
					strings.Fields(fmt.Sprint(version)), "."), "[]"))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
