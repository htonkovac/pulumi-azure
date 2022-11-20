package main

import (
	"bajica/azure"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	// "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		log.Print("Creating RG!")
		res := deploy_pulumi()
		name, ok := res.Outputs["rg_name"].Value.(string)
		if !ok {
			println("NO GUD")
		}
		c.String(http.StatusOK, "RG_NAME: %v", name)
	})

	r.Run()
}


func deploy_pulumi() auto.UpResult {

	ctx := context.Background()

	projectName := "myproject"
	// we use a simple stack name here, but recommend using auto.FullyQualifiedStackName for maximum specificity.
	stackName := "dev"
	// stackName := auto.FullyQualifiedStackName("myOrgOrUser", projectName, stackName)

	// create or select a stack matching the specified name and project.
	// this will set up a workspace with everything necessary to run our inline program (deployFunc)
	s, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, azure.Create)

	fmt.Printf("Created/Selected stack %q\n", stackName)

	w := s.Workspace()

	fmt.Println("Installing the Azure plugin")

	// for inline source programs, we must manage plugins ourselves
	err = w.InstallPlugin(ctx, "azure-native", "v1.86.0")
	if err != nil {
		fmt.Printf("Failed to install program plugins: %v\n", err)
		os.Exit(1)
	}
	_, err = s.Refresh(ctx)
	if err != nil {
		fmt.Printf("Failed to refresh stack: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Refresh succeeded!")
	fmt.Println("Starting update")

	// wire up our update to stream progress to stdout
	stdoutStreamer := optup.ProgressStreams(os.Stdout)

	// run the update to deploy our s3 website
	res, err := s.Up(ctx, stdoutStreamer)
	if err != nil {
		fmt.Printf("Failed to update stack: %v\n\n", err)
		os.Exit(1)
	}

	fmt.Println("Update succeeded!")

	return res
}