package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bglin/postman-api-client/goman"
)

func main() {
	client := goman.New(os.Getenv("POSTMAN_API_KEY"))

	var myworkspace, err = client.GetWorkspace(os.Getenv("WORKSPACE_ID"))
	if err != nil {
		log.Printf("%v", err)
	}
	var environments = myworkspace.Workspace.Environments

	envStruct, err := client.GetEnvironment(environments[0].ID)
	if err != nil {
		log.Printf("%v", err)
	}
	fmt.Printf("There are %d env values in %s \n", len(envStruct.Environment.Values), envStruct.Environment.Name)
	for _, v := range envStruct.Environment.Values {

		fmt.Printf("Variable: %s, Initial Value: %v \n", v.Key, v.Value)
	}
}
