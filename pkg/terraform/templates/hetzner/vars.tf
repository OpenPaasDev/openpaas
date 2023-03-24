variable "hcloud_token" {
  sensitive = true
  type = string
}

variable "server_groups" {
  type = list
  default = [{{ range $key, $value := .ServerGroups}}
   {
    name = "{{ $key }}"
    instance_type = "{{ $value.InstanceType}}"
    num = {{ $value.Num }}
    is_public = {{ $value.IsPublic }}
    volumes = [{{ range $key, $value := $value.Volumes}}
     {
      name = "{{ $value.Name }}"
      path = "{{ $value.Path }}"
      size = {{ $value.Size }}
     },{{ end }}
    ]
   },{{ end }}
  ]
}


variable "ssh_keys" {
  type = list
  default = [{{ range $key, $value := .CloudProviderConfig.ProviderSettings.ssh_keys}}
   "{{ $value }}",{{ end }}
  ]
}

variable "base_server_name" {
  type = string
  default = "{{.CloudProviderConfig.ProviderSettings.resource_names.base_server_name}}"
}

variable "load_balancer_type" {
  type = string
  default = "{{.CloudProviderConfig.ProviderSettings.load_balancer_type}}"
}

variable "firewall_name" {
  type = string
  default = "{{.CloudProviderConfig.ProviderSettings.resource_names.firewall_name}}"
}

variable "network_name" {
  type = string
  default = "{{.CloudProviderConfig.ProviderSettings.resource_names.network_name}}"
}

variable "allow_ips" {
  type = list
  default = [{{ range $key, $value := .CloudProviderConfig.AllowedIPs}}
   "{{ $value }}",{{ end }}
  ]
}

variable "https_allowed_ips" {
  type = list
  default = [{{ range $key, $value := .CloudProviderConfig.ProviderSettings.https_allowed_ips}}
   "{{ $value }}",{{ end }}
  ]
}

variable "ssl_certificate_ids" {
  type = list
  default = [{{ range $key, $value := .CloudProviderConfig.ProviderSettings.ssl_certificate_ids}}
   {{ $value }},{{ end }}
  ]
}

variable "location"{
  type = string
  default = "{{.CloudProviderConfig.ProviderSettings.location}}"
}