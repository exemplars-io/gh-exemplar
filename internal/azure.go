package internal

import (
	"bytes"
	"fmt"
	"log"

	"github.com/Azure/go-autorest/autorest/azure/cli"
)

func GetProfilePath() {

	path, err := cli.AccessTokensPath()
	if err != nil {
		fmt.Println(fmt.Errorf("could not get token path: %w", err).Error())
	}

	fmt.Println("Path is :", path)

	// Get the tokens
	token, err := cli.GetTokenFromCLI("")
	if err != nil {
		fmt.Println(fmt.Errorf("could not get tokens: %w", err).Error())
	}
	fmt.Printf("%v", token.AccessToken)

	// Execute Command
	var cliCmd = cli.GetAzureCLICommand()
	cliCmd.Args = append(cliCmd.Args, "account", "show", "-o", "json")

	var stderr bytes.Buffer
	cliCmd.Stderr = &stderr

	output, err := cliCmd.Output()
	if err != nil {
		if stderr.Len() > 0 {
			fmt.Errorf("Invoking Azure CLI failed with the following error: %s", stderr.String())
		}

		fmt.Errorf("Invoking Azure CLI failed with the following error: %s", err.Error())
	}

	//fmt.Printf("output of command is :\t %v \n", output)
	//var anyJson map[string]interface{}
	//json.Unmarshal(output, &anyJson)

	res, err := prettyJSON(output)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)

}

// internal functions
// this function will create an SPN for exemplar, this will be used to execute terraform action in repository

func CreateExemplarSpn(spnname string, subscriptionid string) (azureserviceprincipal, error) {

	// Config parameters
	spn := azureserviceprincipal{}
	var sid_to_use = subscriptionid
	var err error
	if sid_to_use == "" {

		// Get the default subscription
		sid_to_use, err = DefaultSubscriptionId()

		fmt.Printf("Subscription id provided is empty, using the default subscription id : %s\n", sid_to_use)

		if err != nil {
			fmt.Errorf("Invoking Azure CLI failed with the following error: %s", err.Error())
		}
	}

	// Execute Command
	var cliCmd = cli.GetAzureCLICommand()
	cliCmd.Args = append(cliCmd.Args, "ad", "sp", "create-for-rbac", "-n", spnname, "--role", "owner")

	var stderr bytes.Buffer
	cliCmd.Stderr = &stderr

	output, err := cliCmd.Output()

	if err != nil {
		if stderr.Len() > 0 {
			fmt.Errorf("Invoking Azure CLI failed with the following error: %s", stderr.String())
		}

		fmt.Errorf("Invoking Azure CLI failed with the following error: %s", err.Error())
	}

	// create spn object based on result
	spn.subscriptionId = sid_to_use

	appId, err := getJsonKeyValue(output, "appId")

	spn.appId = appId

	password, err := getJsonKeyValue(output, "password")

	spn.password = password

	displayName, err := getJsonKeyValue(output, "displayName")

	spn.displayName = displayName

	tenant, err := getJsonKeyValue(output, "tenant")
	spn.tenant = tenant

	return spn, nil

}

// function to get default subscription id
// *************************************************
func DefaultSubscriptionId() (string, error) {

	// Execute Account Show
	var cliCmd = cli.GetAzureCLICommand()
	cliCmd.Args = append(cliCmd.Args, "account", "show")

	var stderr bytes.Buffer
	cliCmd.Stderr = &stderr

	output, err := cliCmd.Output()
	if err != nil {
		if stderr.Len() > 0 {
			fmt.Errorf("Invoking Azure CLI failed with the following error: %s", stderr.String())
		}
	}

	// get the id from json block
	subscriptionid, err := getJsonKeyValue(output, "id")

	if err == nil {
		return subscriptionid, nil
	}

	return "", nil

}

// *****************************************************

// function to get azure subscription name
// *****************************************************
func GetSubscriptionName(subscriptionid string) (string, error) {

	// Execute Account Show
	var cliCmd = cli.GetAzureCLICommand()
	cliCmd.Args = append(cliCmd.Args, "account", "show", "-s", subscriptionid)

	var stderr bytes.Buffer
	cliCmd.Stderr = &stderr

	output, err := cliCmd.Output()
	if err != nil {
		if stderr.Len() > 0 {
			fmt.Errorf("Invoking Azure CLI failed with the following error: %s", stderr.String())
		}
	}

	// get the id from json block
	name, err := getJsonKeyValue(output, "name")

	if err == nil {
		return name, nil
	}

	return "", nil
}

// ******************************************************
