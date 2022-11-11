/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/exemplars-io/gh-exemplar/internal"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// check if or and and keywords are provided for execution
		flagTopic, _ := cmd.Flags().GetBool("topic")
		keywords, _ := cmd.Flags().GetString("keywords")

		var profile = internal.GetUserProfile(cmd, args)

		internal.ListExemplars(keywords, flagTopic, profile)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("keywords", "k", "", "keywords to filter exemplars list, use space for multiple keywords by defualt keywords will be searched in name,description and readme files")
	listCmd.MarkFlagRequired("keywords")

	listCmd.Flags().BoolP("topic", "t", false, "search keywoprds in topics")
	//listCmd.Flags().StringP("or", "o", "", "match any keywwords specified")

}
