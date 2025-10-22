resource "aws_subnet" "this" {
  vpc_id                  = var.vpc_id
  cidr_block              = var.public_subnet_cidr
  availability_zone       = var.az
  map_public_ip_on_launch = true
  tags = merge(var.tags, {
    Name = "${var.name}-subnet"
  })
}

