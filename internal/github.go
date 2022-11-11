/*
Copyright © 2022 NAME HERE satya.tanwar@gmail.com

*/
package internal

import (
	"context"
	crypto_rand "crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/cli/go-gh"
	"github.com/google/go-github/v43/github"
	"golang.org/x/crypto/nacl/box"
)

/// external functions

// function get list of exemplars based on keywords and flags
func ListExemplars(Keywords string, topic bool, profile userprofile) {

	var query = ""
	ctx, client, err := githubClient(profile)

	CheckIfError(err)

	// list public repositories for org "github"
	opts := &github.SearchOptions{Sort: "forks", Order: "desc", ListOptions: github.ListOptions{PerPage: 10}}

	// "key1 key2" will fetch matching repo with any keywords and idicate this as OR
	// key1 key2 will fetch matching repo with all keywords and indicate this as AND
	// var query = `org:exemplars-io is:public topic:"k8s elk","demo vemo"`
	// up example will get all repos in org : exemplarsi-io , that are publicably visible and have topics: (k8s or elk) and (demo or vemo)
	if topic {
		query = `org:exemplars-io is:public topic:"` + Keywords + `"`
	} else {
		query = `org:exemplars-io is:public ` + Keywords + ` in:name,description,readme `
	}

	// var query = `org:exemplars-io is:public ` + Keywords + ` in:name,description,readme `
	result, _, err := client.Search.Repositories(ctx, query, opts)

	CheckIfError(err)

	Info("Total: %d exemplars found\n", *result.Total)

	CheckIfError(err)

	for i, repo := range result.Repositories {
		fmt.Printf("%v. %v --- [%v] \n", i+1, repo.GetName(), repo.GetDescription())
	}
}

/// function to prepare exemplar for execution
// 1. Clone repo from exemplar template to user's github profile
func PrepareExemplar(name string, profile userprofile) (*github.Repository, error) {

	// get the existing repo if alreay present
	existingrepo, err := getRepo(name, profile)

	CheckIfError(err)

	// get git hub client
	ctx, client, err := githubClient(profile)

	CheckIfError(err)

	if existingrepo != nil {
		// Repo already exist.. prompt if want to delete and re-create
		Info("Repo exist for exemplar : %v at %v", name, existingrepo.GetCloneURL())
		fmt.Printf("Delete existing repo and create again (y/n): ")
		var input string
		fmt.Scan(&input)
		// Check input
		if input == "y" {
			Info("Deleting the repo now...")
			_, err := client.Repositories.Delete(ctx, profile.owner, name)
			CheckIfError(err)
			Info("Repo : %v deleted succesfully", existingrepo.GetCloneURL())

		} else {
			// return existing repo
			return existingrepo, err
		}

	}
	Info("Create the repo for exemplar now...")
	// Create repo from exemplar
	request := &github.TemplateRepoRequest{
		Name:  &name,
		Owner: &profile.owner,
	}
	// build repo url
	//repo_url := e_url + name
	result, _, err := client.Repositories.CreateFromTemplate(ctx, e_owner, name, request)
	CheckIfError(err)
	Info("Created repo : %v for exemplar - %v", result.GetCloneURL(), name)
	existingrepo = result

	return existingrepo, err

}

// function to execute exemplar

// 1. This will run the complete exemplar
// - Read the Exemplar config file
// - Validate the Exemplar
// - Run IaC Action and wait for succesfull execution
// - Run CI Action to build the CNAB
// - Run Test Action to
func Executexemplar(name string, profile userprofile) {

}

func getRepo(name string, profile userprofile) (*github.Repository, error) {

	// get the git client
	ctx, client, err := githubClient(profile)

	CheckIfError(err)
	var query = ""
	opts := &github.SearchOptions{}
	if profile.organization != "" {
		query = `org:` + profile.organization + ` is:public ` + name + ` in:name`
	} else {
		query = `user:` + profile.owner + ` is:public ` + name + ` in:name`
	}
	result, _, err := client.Search.Repositories(ctx, query, opts)
	if len(result.Repositories) == 1 {
		return result.Repositories[0], err
	}
	return nil, err

}

/// func to add or update repo secrets that can be used by actions
func AddRepoSecret(owner string, repo string, sname string, svalue string, profile userprofile) {

	ctx, client, err := githubClient(profile)
	if err != nil {
		log.Fatal(err)
	}

	// TODO Remove hard coded owner name
	if err := addRepoSecret(ctx, client, owner, repo, sname, svalue); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Added secret %q to the repo %v\n", sname, repo)
}

/// internal methods for github integration

/**
func githubClient(profile userprofile) (context.Context, *github.Client, error) {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: profile.githubtoken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return ctx, client, nil
}

**/

func githubClient(profile userprofile) (context.Context, *github.Client, error) {

	ctx := context.Background()
	//ts := oauth2.StaticTokenSource(
	//		&oauth2.Token{AccessToken: profile.githubtoken},
	//	)
	tc, err := gh.HTTPClient(nil)
	if err != nil {
		log.Fatal(err)
	}

	client := github.NewClient(tc)
	return ctx, client, nil
}

func addRepoSecret(ctx context.Context, client *github.Client, owner string, repo, secretName string, secretValue string) error {
	publicKey, _, err := client.Actions.GetRepoPublicKey(ctx, owner, repo)
	if err != nil {
		return err
	}

	encryptedSecret, err := encryptSecretWithPublicKey(publicKey, secretName, secretValue)
	if err != nil {
		return err
	}

	if _, err := client.Actions.CreateOrUpdateRepoSecret(ctx, owner, repo, encryptedSecret); err != nil {
		return fmt.Errorf("Actions.CreateOrUpdateRepoSecret returned error: %v", err)
	}

	return nil
}

// func to encrypt secret value with public key
func encryptSecretWithPublicKey(publicKey *github.PublicKey, secretName string, secretValue string) (*github.EncryptedSecret, error) {

	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey.GetKey())
	if err != nil {
		return nil, fmt.Errorf("base64.StdEncoding.DecodeString was unable to decode public key: %v", err)
	}

	var boxKey [32]byte
	copy(boxKey[:], decodedPublicKey)
	secretBytes := []byte(secretValue)
	encryptedBytes, err := box.SealAnonymous([]byte{}, secretBytes, &boxKey, crypto_rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("box.SealAnonymous failed with error %w", err)
	}

	encryptedString := base64.StdEncoding.EncodeToString(encryptedBytes)
	keyID := publicKey.GetKeyID()
	encryptedSecret := &github.EncryptedSecret{
		Name:           secretName,
		KeyID:          keyID,
		EncryptedValue: encryptedString,
	}
	return encryptedSecret, nil
}

// func to update secrets for azure
func AddUpdateAzureSecretsToGithubRepo(owner string, repo string, spn azureserviceprincipal, profile userprofile) {

	// Add secret values required for azure provider
	AddRepoSecret(owner, repo, "AZURE_AD_CLIENT_ID", spn.appId, profile)
	AddRepoSecret(owner, repo, "AZURE_AD_CLIENT_SECRET", spn.password, profile)
	AddRepoSecret(owner, repo, "AZURE_AD_TENANT_ID", spn.tenant, profile)
	AddRepoSecret(owner, repo, "AZURE_SUBSCRIPTION_ID", spn.subscriptionId, profile)

}
