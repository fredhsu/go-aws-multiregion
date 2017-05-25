package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// GetEC2PublicIP takes the physical ID of an EC2 instance and a session
// returning the public IP of that instance

func GetEC2PublicIP(ec2id *string, sess client.ConfigProvider) *string {
	if ec2id == nil {
		fmt.Println("Invalid ec2, cannot generate IP")
		return nil
	}

	ec2svc := ec2.New(sess)
	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(*ec2id),
		},
	}
	resp, err := ec2svc.DescribeInstances(params)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return resp.Reservations[0].Instances[0].PublicIpAddress
}

// GetStackEC2ID takes the stack name and a session then returns the
// physical id of the Ec2 instance configured

func GetStackEC2ID(stack string, sess client.ConfigProvider) *string {
	// TODO Verify if the stack exists
	svc := cloudformation.New(sess)
	//params := &cloudformation.DescribeStacksInput{
	//	NextToken: aws.String("Next"),
	//	StackName: aws.String(stack),
	//}
	// resp, err := svc.DescribeStacks(params)

	//	if err != nil {
	//		fmt.Println(err.Error())
	//		return nil
	//	}
	// outputs := resp.Stacks[0].Outputs
	//	fmt.Println(outputs)
	//	for _, o := range outputs {
	//		if *o.OutputKey == "Ec2IP" {
	//			fmt.Printf("ip: %s\n", *o.OutputValue)
	//		}
	//	}
	params2 := &cloudformation.DescribeStackResourcesInput{
		LogicalResourceId: aws.String("Ec2Instance"),
		StackName:         aws.String(stack),
	}
	resp2, err := svc.DescribeStackResources(params2)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return resp2.StackResources[0].PhysicalResourceId

}

/*
getCTConfig takes an ip address and generates a connectivity monitor config for it, priting to stdout:
	monitor connectivity
	   host aws-east2
		  ip 52.14.135.105
		  url http://52.14.135.105
*/
func genCTConfig(region string, ip string) {
	fmt.Println("monitor connectivity")
	fmt.Println("host " + region)
	fmt.Println("ip " + ip)
	fmt.Println("url http://" + ip)
}

func main() {
	// TODO get regions from config file
	regions := []string{"us-west-2", "eu-west-1", "ap-southeast-2", "ap-northeast-1"}
	//regions := []string{"us-west-2"}
	// TODO use goroutines
	for _, re := range regions {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(re)})
		if err != nil {
			fmt.Println(err)
		}

		stack := "CT-" + re
		ec2id := GetStackEC2ID(stack, sess)
		// Verify if the ec2 instance exists - ec2id is nil if didn't exist
		publicip := *GetEC2PublicIP(ec2id, sess)
		// TODO move this to another app
		genCTConfig(re, publicip)
	}
}
