/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/exemplars-io/gh-exemplar/internal"

	"github.com/spf13/cobra"
)

// prepareCmd represents the prepare command
var forkCmd = &cobra.Command{
	Use:   "fork",
	Short: "Command to fork an exemplar for customization",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		name, _ := cmd.Flags().GetString("name")

		var profile = internal.GetUserProfile(cmd, args)

		internal.PrepareExemplar(name, profile)
	},
}

func init() {
	rootCmd.AddCommand(forkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// prepareCmd.PersistentFlags().String("foo", "", "A help for foo")

	forkCmd.Flags().StringP("name", "n", "", "name of exemplar to prepare for execution")
	forkCmd.MarkFlagRequired("name")
}
