package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestSubnetModule(t *testing.T) {
    t.Parallel()

    // First, create a VPC (dependency)
    vpcOptions := &terraform.Options{
        TerraformDir: "../modules/vpc",
        Vars: map[string]interface{}{
            "name":     "test-vpc",
            "vpc_cidr": "10.0.0.0/16",
            "tags": map[string]string{
                "Environment": "dev",
                "Project":     "opentofu-demo",
            },
        },
    }

    defer terraform.Destroy(t, vpcOptions)
    terraform.InitAndApply(t, vpcOptions)
    vpcID := terraform.Output(t, vpcOptions, "vpc_id")

    // Now test the Subnet
    terraformOptions := &terraform.Options{
        TerraformDir: "../modules/subnet",
        Vars: map[string]interface{}{
            "name":               "test-subnet",
            "vpc_id":             vpcID,
            "public_subnet_cidr": "10.0.1.0/24",
            "az":                 "ap-south-1a",
            "tags": map[string]string{
                "Environment": "dev",
                "Project":     "opentofu-demo",
            },
        },
    }

    defer terraform.Destroy(t, terraformOptions)
    terraform.InitAndApply(t, terraformOptions)

    subnetID := terraform.Output(t, terraformOptions, "subnet_id")
    assert.NotEmpty(t, subnetID)
}

