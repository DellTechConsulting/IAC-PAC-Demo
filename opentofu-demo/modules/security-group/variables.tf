variable "name" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "my_ip_cidr" {
  type = string
}

variable "tags" {
  type    = map(string)
  default = {}
}

