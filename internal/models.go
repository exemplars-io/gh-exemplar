package internal

// Struct to represent Azure Service Principal object
type azureserviceprincipal struct {
	appId          string
	tenant         string
	password       string
	displayName    string
	subscriptionId string
}

// model for user profile loaded by config file
type userprofile struct {
	owner        string
	organization string
	githubtoken  string
}

const e_owner = "exemplars-io"
