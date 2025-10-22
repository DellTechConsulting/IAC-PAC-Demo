terraform {
  required_providers {
    aws = {
      source  = "registry.opentofu.org/hashicorp/aws"
      version = ">= 5.0.0"
    }
  }
}

provider "aws" {
  region = "ap-south-1"
}

module "vpc" {
  source   = "../../modules/vpc"
  name     = "dev-vpc"
  vpc_cidr = "10.0.0.0/16"
  tags = {
    Environment = "dev"
    Project     = "opentofu-demo"
  }
}

module "subnet" {
  source             = "../../modules/subnet"
  name               = "dev-subnet"
  vpc_id             = module.vpc.vpc_id
  public_subnet_cidr = "10.0.1.0/24"
  az                 = "ap-south-1a"
  tags = {
    Environment = "dev"
    Project     = "opentofu-demo"
  }
}

module "security_group" {
  source     = "../../modules/security-group"
  name       = "dev-sg"
  vpc_id     = module.vpc.vpc_id
  my_ip_cidr = "0.0.0.0/0"
  tags = {
    Environment = "dev"
    Project     = "opentofu-demo"
  }
}

module "ec2" {
  source        = "../../modules/ec2"
  name          = "dev-ec2"
  ami_id        = "ami-0dee22c13ea7a9a67"
  instance_type = "t3.micro"
  subnet_id     = module.subnet.subnet_id
  sg_id         = module.security_group.sg_id
  tags = {
    Environment = "dev"
    Project     = "opentofu-demo"
  }
}

