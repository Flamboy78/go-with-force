package main

import (
	"fmt"
	"klnef/go-with-force/internal/credentials"
	"klnef/go-with-force/internal/soql"
	"net/http"
)

func main() {
	creds, err := credentials.GeneratePasswordCredentials(credentials.PasswordCredentials{
		Username:     "testUser",                     // Change it with the username of your Salesfroce User
		Password:     "testPassword",                 // Change it with the password of your Salesfroce User
		ClientId:     "12424214",                     // The Client Id (Consumer Key) of your Connected App
		ClientSecret: "233225324",                    // The Client Secret (Consumer Secret) of your Connected App
		Url:          "https://login.salesforce.com", // Url of your Salesforce Org (Can leave as it is if you are not using a sandbox)
	},
	)

	if err != nil {
		fmt.Printf("error %v\n", err)
		return
	}
	setup := credentials.Setup{
		PasswordCredentials: *creds,
		Protocol:            http.DefaultClient,
	}

	session, err := credentials.Auth(setup)

	if err != nil {
		fmt.Printf("error %v\n", err)
		return
	}

	fmt.Printf("Success %v\n", session.PassCred)

	resource, err := soql.NewResource(session)
	if err != nil {
		fmt.Printf("Error creating resource: %v\n", err)
		return
	}

	result, err := resource.Query("SELECT Id, Name FROM Account") // Sample query for testing purposes
	if err != nil {
		fmt.Printf("SOQL Query Error %s\n", err.Error())
		return
	}
	fmt.Println("SOQL Query")
	fmt.Println("---------------------")
	fmt.Printf("Done: %v\n", result)

}
