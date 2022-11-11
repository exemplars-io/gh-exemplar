package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/spf13/cobra"
)

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// function to get a field from json block
// uses the below library for parser
// https://github.com/buger/jsonparser
func getJsonKeyValue(data []byte, keys ...string) (string, error) {

	val, err := jsonparser.GetString(data, keys...)
	if err != nil {
		return "", err
	}
	return val, nil
}

/// function to show a json string into pretty/indented format
func prettyJSON(input []byte) (string, error) {

	var prettyJSON bytes.Buffer

	if err := json.Indent(&prettyJSON, input, "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

/// function load user profile based on config
func GetUserProfile(cmd *cobra.Command, args []string) userprofile {

	// get required flag values
	f_githubtoken, _ := cmd.Flags().GetString("gittoken")
	f_owner, _ := cmd.Flags().GetString("gitowner")
	f_organization, _ := cmd.Flags().GetString("gitorg")

	var profile = userprofile{
		githubtoken:  f_githubtoken,
		owner:        f_owner,
		organization: f_organization,
	}

	return profile
}
func GetNilUserProfile() userprofile {

	var profile = userprofile{
		githubtoken:  "nil",
		owner:        "nil",
		organization: "nil",
	}

	return profile
}
