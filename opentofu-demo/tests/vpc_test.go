package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestVPCModule(t *testing.T) {
    terraformOptions := &terraform.Options{
        TerraformDir: "../modules/vpc",
        Vars: map[string]interface{}{
            "name":     "test-vpc",
            "vpc_cidr": "10.0.0.0/16",
            "tags": map[string]string{
                "Environment": "dev",
                "Module":      "vpc",
            },
        },
    }

    // Always attempt cleanup — even on failure
    defer func() {
        terraform.DestroyE(t, terraformOptions)
    }()

    terraform.InitAndApply(t, terraformOptions)

    vpcID := terraform.Output(t, terraformOptions, "vpc_id")
    assert.NotEmpty(t, vpcID, "VPC ID should not be empty")
}

