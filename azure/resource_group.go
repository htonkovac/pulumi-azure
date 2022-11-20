package azure

import (
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/core"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"log"
)

func Create(ctx *pulumi.Context) error {
	log.Print("HELLO")
	rg, err := core.NewResourceGroup(ctx, "example", &core.ResourceGroupArgs{
		Location: pulumi.String("West Europe"),
	})

	ctx.Export("rg_name",rg.ToResourceGroupOutput().Name())
	if err != nil {
		return err
	}
	return nil
}
 