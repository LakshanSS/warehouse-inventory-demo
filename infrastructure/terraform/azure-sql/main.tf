# Azure SQL Database Terraform Configuration
# Provisions SQL Server and Database for the warehouse inventory demo

terraform {
  required_version = ">= 1.5.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }

  # Backend configuration - override via -backend-config flags
  backend "azurerm" {
    resource_group_name  = "terraform-state-rg"
    storage_account_name = "tfstatestore"
    container_name       = "tfstate"
    key                  = "database.tfstate"
  }
}

provider "azurerm" {
  features {}
}

# -----------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------

variable "resource_group_name" {
  description = "Azure Resource Group name"
  type        = string
}

variable "location" {
  description = "Azure region for resources"
  type        = string
  default     = "eastus"
}

variable "server_name" {
  description = "SQL Server name (must be globally unique)"
  type        = string
}

variable "database_name" {
  description = "Database name"
  type        = string
}

variable "database_sku" {
  description = "Database SKU/pricing tier"
  type        = string
  default     = "Basic"

  validation {
    condition     = contains(["Basic", "Standard_S0", "Standard_S1", "Standard_S2", "Premium_P1"], var.database_sku)
    error_message = "Invalid SKU. Must be one of: Basic, Standard_S0, Standard_S1, Standard_S2, Premium_P1"
  }
}

variable "admin_username" {
  description = "SQL Server admin username"
  type        = string
  sensitive   = true
}

variable "admin_password" {
  description = "SQL Server admin password"
  type        = string
  sensitive   = true

  validation {
    condition     = length(var.admin_password) >= 12
    error_message = "Password must be at least 12 characters long"
  }
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default = {
    managed_by = "openchoreo-workflow"
    project    = "warehouse-inventory-demo"
  }
}

# -----------------------------------------------------------------------------
# Resources
# -----------------------------------------------------------------------------

# SQL Server
resource "azurerm_mssql_server" "main" {
  name                         = var.server_name
  resource_group_name          = var.resource_group_name
  location                     = var.location
  version                      = "12.0"
  administrator_login          = var.admin_username
  administrator_login_password = var.admin_password

  minimum_tls_version = "1.2"

  tags = merge(var.tags, {
    resource_type = "sql-server"
  })
}

# SQL Database
resource "azurerm_mssql_database" "main" {
  name         = var.database_name
  server_id    = azurerm_mssql_server.main.id
  sku_name     = var.database_sku
  collation    = "SQL_Latin1_General_CP1_CI_AS"
  max_size_gb  = var.database_sku == "Basic" ? 2 : 10

  # Prevent accidental deletion
  lifecycle {
    prevent_destroy = false
  }

  tags = merge(var.tags, {
    resource_type = "sql-database"
  })
}

# Firewall rule - Allow Azure services
resource "azurerm_mssql_firewall_rule" "allow_azure_services" {
  name             = "AllowAzureServices"
  server_id        = azurerm_mssql_server.main.id
  start_ip_address = "0.0.0.0"
  end_ip_address   = "0.0.0.0"
}

# -----------------------------------------------------------------------------
# Outputs
# -----------------------------------------------------------------------------

output "server_name" {
  description = "SQL Server name"
  value       = azurerm_mssql_server.main.name
}

output "server_fqdn" {
  description = "SQL Server fully qualified domain name"
  value       = azurerm_mssql_server.main.fully_qualified_domain_name
}

output "database_name" {
  description = "Database name"
  value       = azurerm_mssql_database.main.name
}

output "database_id" {
  description = "Database resource ID"
  value       = azurerm_mssql_database.main.id
}

output "connection_string" {
  description = "ADO.NET connection string"
  value       = "Server=tcp:${azurerm_mssql_server.main.fully_qualified_domain_name},1433;Database=${azurerm_mssql_database.main.name};User ID=${var.admin_username};Password=${var.admin_password};Encrypt=true;TrustServerCertificate=false;Connection Timeout=30;"
  sensitive   = true
}

output "jdbc_connection_string" {
  description = "JDBC connection string"
  value       = "jdbc:sqlserver://${azurerm_mssql_server.main.fully_qualified_domain_name}:1433;database=${azurerm_mssql_database.main.name};user=${var.admin_username};password=${var.admin_password};encrypt=true;trustServerCertificate=false;loginTimeout=30;"
  sensitive   = true
}
