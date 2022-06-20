package main

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AppProps struct {
	name       string
	StackProps awscdk.StackProps
}

func NewCDKStack(scope constructs.Construct, appName string, props *AppProps) awscdk.Stack {

	stack := awscdk.NewStack(scope, &appName, &props.StackProps)
	bundlingOptions := &awscdklambdagoalpha.BundlingOptions{
		GoBuildFlags: &[]*string{jsii.String(fmt.Sprintf(`-ldflags "%s"`, "-s -w"))},
	}

	quadraticLambda := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("Quadratic Get Lambda"), &awscdklambdagoalpha.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../api/handlers/maths/quadratic/get"),
		Bundling:     bundlingOptions,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(10)),
		Tracing:      awslambda.Tracing_ACTIVE,
		MemorySize:   jsii.Number(1024),
		Environment: &map[string]*string{
			"OFFSET": jsii.String("BLAH"),
		},
	})

	quadraticLambda.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_AWS_IAM,
	})

	addLambda := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("Add Get Lambda"), &awscdklambdagoalpha.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../api/handlers/maths/add/get"),
		Bundling:     bundlingOptions,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(10)),
		Tracing:      awslambda.Tracing_ACTIVE,
		MemorySize:   jsii.Number(1024),
	})

	// Not Found handler.
	notFoundLambda := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("notFoundHandler"), &awscdklambdagoalpha.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../api/handlers/notfound"),
		Bundling:     bundlingOptions,
	})

	addLambda.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_AWS_IAM,
	})

	api := awsapigateway.NewLambdaRestApi(stack, jsii.String("example-api"), &awsapigateway.LambdaRestApiProps{
		Handler: notFoundLambda,
		Proxy:   jsii.Bool(false),
	})
	apiResourceOpts := &awsapigateway.ResourceOptions{}
	apiLambdaOpts := &awsapigateway.LambdaIntegrationOptions{}
	iamAuthMethodOps := &awsapigateway.MethodOptions{
		AuthorizationType: awsapigateway.AuthorizationType_IAM,
	}

	maths := api.Root().AddResource(jsii.String("maths"), apiResourceOpts)
	maths.AddResource(jsii.String("add"), apiResourceOpts).
		AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(addLambda, apiLambdaOpts), iamAuthMethodOps)
	maths.AddResource(jsii.String("quadratic"), apiResourceOpts).
		AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(addLambda, apiLambdaOpts), iamAuthMethodOps)

	awscdk.NewCfnOutput(stack, jsii.String("QuadraticLambdaARN"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("QuadraticLambdaARN"),
		Value:      quadraticLambda.FunctionArn(),
	})
	awscdk.NewCfnOutput(stack, jsii.String("AddLambdaARN"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("AddLambdaARN"),
		Value:      addLambda.FunctionArn(),
	})

	return stack
}

func main() {

	app := awscdk.NewApp(nil)

	props := &AppProps{
		name: "example-api",
	}

	NewCDKStack(app, props.name, props)

	app.Synth(nil)

}
