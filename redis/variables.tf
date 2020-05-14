
variable "client_id" {}
variable "client_secret" {}

variable rg_name {
    default = "mredisrg"
}

variable location {
    default = "westeurope"
}

variable redis_name {
    default = "mredis"
}
