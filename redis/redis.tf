
resource "azurerm_resource_group" "redisrg" {
    name     = var.rg_name
    location = var.location
}

resource "azurerm_redis_cache" "main" {
    name                = var.redis_name
    location            = azurerm_resource_group.redisrg.location
    resource_group_name = azurerm_resource_group.redisrg.name
    capacity            = 0
    family              = "C"
    sku_name            = "Basic"
    enable_non_ssl_port = false
    minimum_tls_version = "1.2"

    redis_configuration {
    }
}

data "azurerm_redis_cache" "dataredis" {
  name                = azurerm_redis_cache.main.name
  resource_group_name = azurerm_resource_group.redisrg.name
}

output "primary_access_key" {
  value = data.azurerm_redis_cache.dataredis.primary_access_key
}

output "hostname" {
  value = data.azurerm_redis_cache.dataredis.hostname
}

output "ssl_port" {
  value = data.azurerm_redis_cache.dataredis.ssl_port
}