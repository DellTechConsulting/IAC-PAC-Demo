package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestSecurityGroupModule(t *testing.T) {
    t.Parallel()

    // Create VPC dependency
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

    // Test SG module
    terraformOptions := &terraform.Options{
        TerraformDir: "../modules/security-group",
        Vars: map[string]interface{}{
            "name":       "test-sg",
            "vpc_id":     vpcID,
            "my_ip_cidr": "1.2.3.4/32",
            "tags": map[string]string{
                "Environment": "dev",
                "Project":     "opentofu-demo",
            },
        },
    }

    defer terraform.Destroy(t, terraformOptions)
    terraform.InitAndApply(t, terraformOptions)

    sgID := terraform.Output(t, terraformOptions, "sg_id")
    assert.NotEmpty(t, sgID)
}

