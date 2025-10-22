package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestSubnetModule(t *testing.T) {
    // Create a VPC first (required for subnet)
    vpcOptions := &terraform.Options{
        TerraformDir: "../modules/vpc",
        Vars: map[string]interface{}{
            "name":     "test-vpc-for-subnet",
            "vpc_cidr": "10.1.0.0/16",
        },
    }

    defer func() {
        terraform.DestroyE(t, vpcOptions)
    }()
    terraform.InitAndApply(t, vpcOptions)

    vpcID := terraform.Output(t, vpcOptions, "vpc_id")
    assert.NotEmpty(t, vpcID)

    // Now create a subnet
    subnetOptions := &terraform.Options{
        TerraformDir: "../modules/subnet",
        Vars: map[string]interface{}{
            "name":               "test-subnet",
            "vpc_id":             vpcID,
            "public_subnet_cidr": "10.1.1.0/24",
            "az":                 "ap-south-1a",
        },
    }

    defer func() {
        terraform.DestroyE(t, subnetOptions)
    }()

    terraform.InitAndApply(t, subnetOptions)

    subnetID := terraform.Output(t, subnetOptions, "subnet_id")
    assert.NotEmpty(t, subnetID, "Subnet ID should not be empty")
}

