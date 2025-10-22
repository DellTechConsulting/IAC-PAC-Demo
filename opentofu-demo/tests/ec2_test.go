package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestEC2Module(t *testing.T) {
    // Create a VPC
    vpcOptions := &terraform.Options{
        TerraformDir: "../modules/vpc",
        Vars: map[string]interface{}{
            "name":     "test-vpc-for-ec2",
            "vpc_cidr": "10.3.0.0/16",
        },
    }

    defer func() {
        terraform.DestroyE(t, vpcOptions)
    }()
    terraform.InitAndApply(t, vpcOptions)
    vpcID := terraform.Output(t, vpcOptions, "vpc_id")

    // Create a subnet
    subnetOptions := &terraform.Options{
        TerraformDir: "../modules/subnet",
        Vars: map[string]interface{}{
            "name":               "test-subnet-ec2",
            "vpc_id":             vpcID,
            "public_subnet_cidr": "10.3.1.0/24",
            "az":                 "ap-south-1a",
        },
    }

    defer func() {
        terraform.DestroyE(t, subnetOptions)
    }()
    terraform.InitAndApply(t, subnetOptions)
    subnetID := terraform.Output(t, subnetOptions, "subnet_id")

    // Create security group
    sgOptions := &terraform.Options{
        TerraformDir: "../modules/security-group",
        Vars: map[string]interface{}{
            "name":       "test-sg-ec2",
            "vpc_id":     vpcID,
            "my_ip_cidr": "1.2.3.4/32",
        },
    }

    defer func() {
        terraform.DestroyE(t, sgOptions)
    }()
    terraform.InitAndApply(t, sgOptions)
    sgID := terraform.Output(t, sgOptions, "sg_id")

    // Launch EC2
    ec2Options := &terraform.Options{
        TerraformDir: "../modules/ec2",
        Vars: map[string]interface{}{
            "name":          "test-ec2",
            "ami_id":        "ami-0c2b8ca1dad447f8a", // ap-south-1 Amazon Linux 2
            "instance_type": "t3.micro",
            "subnet_id":     subnetID,
            "sg_id":         sgID,
        },
    }

    defer func() {
        terraform.DestroyE(t, ec2Options)
    }()
    terraform.InitAndApply(t, ec2Options)
    instanceID := terraform.Output(t, ec2Options, "instance_id")

    assert.NotEmpty(t, instanceID, "Instance ID should not be empty")
}

