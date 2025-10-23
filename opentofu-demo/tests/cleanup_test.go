package test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// TestCleanupAll ensures all "test" prefixed VPCs, subnets, SGs, and EC2 instances are deleted after tests.
func TestCleanupAll(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		t.Fatalf("Failed to load AWS config: %v", err)
	}

	svc := ec2.NewFromConfig(cfg)
	fmt.Println("🔧 Running forced cleanup of test resources...")

	// 1️⃣ Cleanup EC2 instances first
	instances, err := svc.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err == nil {
		for _, r := range instances.Reservations {
			for _, inst := range r.Instances {
				for _, tag := range inst.Tags {
					if strings.Contains(strings.ToLower(*tag.Value), "test") {
						fmt.Printf("Terminating instance: %s\n", *inst.InstanceId)
						_, _ = svc.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
							InstanceIds: []string{*inst.InstanceId},
						})
					}
				}
			}
		}
	}

	// 2️⃣ Cleanup security groups
	sgs, err := svc.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{})
	if err == nil {
		for _, sg := range sgs.SecurityGroups {
			if strings.Contains(strings.ToLower(*sg.GroupName), "test") {
				fmt.Printf("Deleting security group: %s\n", *sg.GroupId)
				_, _ = svc.DeleteSecurityGroup(ctx, &ec2.DeleteSecurityGroupInput{
					GroupId: sg.GroupId,
				})
			}
		}
	}

	// 3️⃣ Cleanup subnets
	subnets, err := svc.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{})
	if err == nil {
		for _, sn := range subnets.Subnets {
			for _, tag := range sn.Tags {
				if strings.Contains(strings.ToLower(*tag.Value), "test") {
					fmt.Printf("Deleting subnet: %s\n", *sn.SubnetId)
					_, _ = svc.DeleteSubnet(ctx, &ec2.DeleteSubnetInput{
						SubnetId: sn.SubnetId,
					})
				}
			}
		}
	}

	// 4️⃣ Cleanup VPCs last
	vpcs, err := svc.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{})
	if err == nil {
		for _, v := range vpcs.Vpcs {
			for _, tag := range v.Tags {
				if strings.Contains(strings.ToLower(*tag.Value), "test") {
					fmt.Printf("Deleting VPC: %s\n", *v.VpcId)
					_, _ = svc.DeleteVpc(ctx, &ec2.DeleteVpcInput{
						VpcId: v.VpcId,
					})
				}
			}
		}
	}

	fmt.Println("✅ Cleanup complete.")
}

