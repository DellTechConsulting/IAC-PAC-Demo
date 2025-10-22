package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestSecurityGroupModule(t *testing.T) {
    // Create a VPC first (required for SG)
    vpcOptions := &terraform.Options{
        TerraformDir: "../modules/vpc",
        Vars: map[string]interface{}{
            "name":     "test-vpc-for-sg",
            "vpc_cidr": "10.2.0.0/16",
        },
    }

    defer func() {
        terraform.DestroyE(t, vpcOptions)
    }()
    terraform.InitAndApply(t, vpcOptions)

    vpcID := terraform.Output(t, vpcOptions, "vpc_id")
    assert.NotEmpty(t, vpcID)

    // Create SG
    sgOptions := &terraform.Options{
        TerraformDir: "../modules/security-group",
        Vars: map[string]interface{}{
            "name":       "test-sg",
            "vpc_id":     vpcID,
            "my_ip_cidr": "1.2.3.4/32",
        },
    }

    defer func() {
        terraform.DestroyE(t, sgOptions)
    }()

    terraform.InitAndApply(t, sgOptions)

    sgID := terraform.Output(t, sgOptions, "sg_id")
    assert.NotEmpty(t, sgID, "Security Group ID should not be empty")
}

