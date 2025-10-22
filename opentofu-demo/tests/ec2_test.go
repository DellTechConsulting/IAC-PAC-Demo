package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestEC2Module(t *testing.T) {
    t.Parallel()

    // Create dependencies (VPC → Subnet → SG)
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

    subnetOptions := &terraform.Options{
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
    defer terraform.Destroy(t, subnetOptions)
    terraform.InitAndApply(t, subnetOptions)
    subnetID := terraform.Output(t, subnetOptions, "subnet_id")

    sgOptions := &terraform.Options{
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
    defer terraform.Destroy(t, sgOptions)
    terraform.InitAndApply(t, sgOptions)
    sgID := terraform.Output(t, sgOptions, "sg_id")

    // Now test EC2
    terraformOptions := &terraform.Options{
        TerraformDir: "../modules/ec2",
        Vars: map[string]interface{}{
            "name":          "test-ec2",
            "ami_id":        "ami-0f5ee92e2d63afc18", // ✅ Replace with valid AMI ID for ap-south-1
            "instance_type": "t3.micro",
            "subnet_id":     subnetID,
            "sg_id":         sgID,
            "tags": map[string]string{
                "Environment": "dev",
                "Project":     "opentofu-demo",
            },
        },
    }

    defer terraform.Destroy(t, terraformOptions)
    terraform.InitAndApply(t, terraformOptions)

    instanceID := terraform.Output(t, terraformOptions, "instance_id")
    assert.NotEmpty(t, instanceID)
}

