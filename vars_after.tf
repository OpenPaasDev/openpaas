variable "hcloud_token" {
  sensitive = true
  type = string
}

variable "server_groups" {
  type = list
  default = [
   {
    name = "observability"
    instance_type = "cx21"
    num = 1
    lb_target = false
    image = "ubuntu-22.04"
    volumes = [
    ]
   },
   {
    name = "postgres"
    instance_type = "cx21"
    num = 1
    lb_target = false
    image = "ubuntu-22.04"
    volumes = [
    ]
   },
   {
    name = "web"
    instance_type = "cx21"
    num = 1
    lb_target = true
    image = "ubuntu-22.04"
    volumes = [
    ]
   },
  ]
}

variable "ssh_keys" {
  type = list
  default = [
   "19772877",
   "19772878",
   "19772879",
   "19772880",
   "19772881",
   "19772882",
   "19772883",
  ]
}

variable "base_server_name" {
  type = string
  default = "prod"
}

variable "load_balancer_type" {
  type = string
  default = "lb11"
}

variable "firewall_name" {
  type = string
  default = "prod_firewall"
}

variable "network_name" {
  type = string
  default = "prod_network"
}

variable "allow_ips" {
  type = list
  default = [
   "178.197.214.192/32",
   "188.63.149.213/32",
  ]
}

variable "ssl_certificate_ids" {
  type = list
  default = [
   1263298,
  ]
}

variable "location"{
  type = string
  default = "fsn1"
}
