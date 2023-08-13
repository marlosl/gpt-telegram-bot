package awsdeploy

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Deploy(ctx *pulumi.Context) error {
	fmt.Println("Creating log groups...")
	err := CreateLogGroups(ctx)
	if err != nil {
		fmt.Printf("Can't create the log groups: %v\n", err)
		return err
	}

	fmt.Println("Creating policies...")
	err = CreateLambdaRolePolicy(ctx)
	if err != nil {
		fmt.Printf("Can't create policies: %v\n", err)
		return err
	}

	fmt.Println("Creating tables...")
	err = CreateDynamoDbTable(ctx)
	if err != nil {
		fmt.Printf("Can't create tables: %v\n", err)
		return err
	}

	fmt.Println("Creating queue...")
	err = CreateSendImageQueue(ctx)
	if err != nil {
		fmt.Printf("Can't create queue: %v\n", err)
		return err
	}

	fmt.Println("Creating lambda...")
	err = CreateLambdaFunctions(ctx)
	if err != nil {
		fmt.Printf("Can't create lambda: %v\n", err)
		return err
	}

	fmt.Println("Creating event source mapping...")
	err = CreateLambdaEventSourceMapping(ctx)
	if err != nil {
		fmt.Printf("Can't create event source mapping: %v\n", err)
		return err
	}

	return nil
}
