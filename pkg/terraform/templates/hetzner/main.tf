terraform {
  required_providers {
    hcloud = {
      source  = "hetznercloud/hcloud"
      version = "1.44.1"
    }
  }
  {{if eq .TfState.Backend "s3"}}
  backend "s3" {
    endpoints                   = { s3 = "{{.TfState.Config.endpoint }}" }
    bucket                      = "{{.TfState.Config.bucket }}"
    key                         = "openpaas/terraform.tfstate"
    region                      = "{{.TfState.Config.region }}"
    access_key                  = "{{.TfState.Config.access_key }}"
    secret_key                  = "{{.TfState.Config.secret_key }}"
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    skip_requesting_account_id  = true
    skip_s3_checksum = true
    force_path_style            = true
}
  {{ end }}
}

# Configure the Hetzner Cloud Provider
provider "hcloud" {
  token = var.hcloud_token
}

locals {

  groups = {
    for i, value in var.server_groups :
    value.name => {
      count  = value.num,
      subnet = "${i}", group = i, server_type = value.instance_type, lb_target = value.lb_target,
      volumes = value.volumes
    }
  }

  servers = flatten([
    for name, value in local.groups : [
      for i in range(value.count) : {
        group_name  = name,
        private_ip  = "10.0.${value.subnet}.${i + 2}",
        name        = "${var.base_server_name}-${name}-${i + 1}",
        group       = value.group
        index       = i
        lb_target   = value.lb_target
        server_type = value.server_type
        volumes     = value.volumes
      }
    ]
  ])

  volumes = flatten([
    for i, value in local.servers : [
      for j, vol in value.volumes : {
        name = "${var.base_server_name}-${value.group_name}-${value.index + 1}-${j}-${vol.name}",
        server = "${value.name}",
        path = vol.path,
        size = vol.size
      }
    ]
  ])

  placement_groups = length(var.server_groups)
}

resource "hcloud_network" "private_network" {
  name     = var.network_name
  ip_range = "10.0.0.0/16"
}

resource "hcloud_network_subnet" "network_subnet" {
  for_each     = local.groups
  network_id   = hcloud_network.private_network.id
  type         = "cloud"
  network_zone = "eu-central"
  ip_range     = "10.0.${each.value.subnet}.0/24"
}

resource "hcloud_placement_group" "placement_group" {
  count = local.placement_groups
  name  = "server_placement_spread_group-${count.index}"
  type  = "spread"
}


resource "hcloud_firewall" "network_firewall" {
  name = var.firewall_name
  rule {
    direction  = "in"
    protocol   = "tcp"
    port       = "1-10000"
    source_ips = var.allow_ips
  }

  rule {
    direction = "in"
    protocol  = "icmp"
    source_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }

  rule {
    direction  = "in"
    protocol   = "tcp"
    port       = "any"
    source_ips = flatten([
      for index, node in hcloud_server.server_node : [
        "${node.ipv4_address}/32"
      ]
    ])
  }

  rule {
    direction  = "in"
    protocol   = "udp"
    port       = "any"
    source_ips = flatten([
      for index, node in hcloud_server.server_node : [
        "${node.ipv4_address}/32"
      ]
    ])
  }

  depends_on = [
    hcloud_server.server_node,
  ]
}

resource "hcloud_firewall_attachment" "fw_ref" {
    firewall_id = hcloud_firewall.network_firewall.id
    server_ids  = flatten([
      for index, node in hcloud_server.server_node : [
        node.id
      ]
    ])
}

resource "hcloud_server" "server_node" {
  for_each           = { for entry in local.servers : "${entry.name}" => entry }
  name               = each.value.name
  image              = "ubuntu-22.04"
  server_type        = each.value.server_type
  location           = var.location
  placement_group_id = hcloud_placement_group.placement_group[each.value.group].id

  public_net {
    ipv4_enabled = true
    ipv6_enabled = false
  }
  depends_on = [
    hcloud_network_subnet.network_subnet,
  ]

  labels = {
    "group" = each.value.group_name
  }

  ssh_keys = [for id in var.ssh_keys : id]

  # don't destroy existing machines if some data changes
  lifecycle {
    ignore_changes = [
      ssh_keys, user_data
    ]
  }
}

resource "hcloud_server_network" "network_binding" {
  for_each   = { for entry in local.servers : "${entry.name}" => entry }
  server_id  = hcloud_server.server_node[each.value.name].id
  network_id = hcloud_network.private_network.id
  ip         = each.value.private_ip
}


resource "hcloud_volume" "volumes" {
  for_each = { for entry in local.volumes : "${entry.name}" => entry.size }
  location = var.location
  name     = "${each.key}"
  size     = each.value
  format   = "ext4"
  depends_on = [
    hcloud_server.server_node 
  ]
}

resource "hcloud_volume_attachment" "volumes" {
  for_each = { for index, entry in local.volumes : entry.name => entry.server }
  volume_id = hcloud_volume.volumes[each.key].id
  server_id = hcloud_server.server_node[each.value].id
  automount = true

  depends_on = [
    hcloud_volume.volumes 
  ]
}


resource "hcloud_load_balancer" "lb1" {
  name               = "lb1"
  load_balancer_type = var.load_balancer_type
  # network_zone       =  hcloud_network_subnet.network_subnet["consul"].network_zone
  location = var.location
  depends_on = [
    hcloud_server.server_node,
    hcloud_server_network.network_binding,
    hcloud_network_subnet.network_subnet,
  ]
}

resource "hcloud_load_balancer_network" "srvnetwork" {
  load_balancer_id = hcloud_load_balancer.lb1.id
  network_id       = hcloud_network.private_network.id
  ip               = "10.0.0.254" # max 5 consul servers, so 10.0.0.7 is free
  depends_on = [
    hcloud_network.private_network
  ]
}

resource "hcloud_load_balancer_service" "load_balancer_service" {
  load_balancer_id = hcloud_load_balancer.lb1.id
  protocol         = "https"
  destination_port = 80
  http {
    certificates = var.ssl_certificate_ids
  }
}


# # this is unfortunately necessary, because no amount of `depends_on` on the load_balancer_target will ensure
# # the nodes and networks are ready for load_balancer target attachment, other than waiting
resource "time_sleep" "wait" {
  create_duration = "2s"
  depends_on = [
    hcloud_server.server_node,
    hcloud_server_network.network_binding,
    hcloud_network_subnet.network_subnet,
    hcloud_load_balancer_network.srvnetwork
  ]
}


resource "hcloud_load_balancer_target" "load_balancer_target" {
  for_each         = { for key, val in local.servers : val.index => val.name if val.lb_target == true }
  type             = "server"
  load_balancer_id = hcloud_load_balancer.lb1.id
  server_id        = hcloud_server.server_node[each.value].id
  use_private_ip   = true
  depends_on       = [time_sleep.wait]
}


output "servers" {
  value = flatten([
    for index, node in hcloud_server.server_node : [
      for server in local.servers :
      { host       = "${node.ipv4_address}",
        host_name  = "${node.name}",
        private_ip = "${server.private_ip}",
        server_id  = node.id
        group      = node.labels["group"]
      } if server.name == node.name
    ]
  ])
}

output "volumes" {
 value = flatten([
    for index, attachment in hcloud_volume_attachment.volumes : [
      for vol in local.volumes :
      {mount = "/mnt/HC_Volume_${attachment.volume_id}", 
      path = vol.path,
      name = vol.name,
      server_id = attachment.server_id} if hcloud_volume.volumes[index].name == vol.name
    ] 
  ])
}

