/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/exemplars-io/gh-exemplar/internal"

	"github.com/spf13/cobra"
)

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Command to exeucte an exemplar",
	Long:  `This command execute an exemplar whose name is supplied as part of the command argument. This should prA longer description that spans multiple lines and likely contains examples`,
	Run: func(cmd *cobra.Command, args []string) {
		// call execute command and pass all arguments
		exemplars_execute(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// executeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// executeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	executeCmd.Flags().StringP("name", "n", "", "Name of the exemplar to execute")
	executeCmd.Flags().StringP("location", "l", "", "Location of the exemplar to execute")

	executeCmd.MarkFlagsMutuallyExclusive("name", "location")

	//executeCmd.Flags().StringP("location", "l", "", "Location of the exemplar to execute")

}

func exemplars_execute(cmd *cobra.Command, args []string) {

	// check if name and path are provided for execution
	exemplar_name, _ := cmd.Flags().GetString("name")
	//exemplar_location, _ := cmd.Flags().GetString("location")

	var profile = internal.GetUserProfile(cmd, args)

	fmt.Println("executing exemplar : " + exemplar_name)

	// Create/Update Exemplar SPN
	spn, err := internal.CreateExemplarSpn("exemplars-spn", "")
	if err == nil {
		// Update Github repo with SPN secrets
		internal.AddUpdateAzureSecretsToGithubRepo("exemplars-io", "exemplar-demo", spn, profile)

	} else {

		fmt.Errorf("error: %s", err.Error())

	}

}
