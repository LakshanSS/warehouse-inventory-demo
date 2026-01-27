# Variables file for Azure SQL Database provisioning
# These can be overridden via -var flags or .tfvars files

# This file documents all available variables with their defaults
# The actual variable definitions are in main.tf to keep things simple

# Example usage:
#
# terraform plan \
#   -var="resource_group_name=my-rg" \
#   -var="server_name=my-sql-server" \
#   -var="database_name=my-db" \
#   -var="admin_username=sqladmin" \
#   -var="admin_password=SecureP@ssw0rd123"
#
# Or with a .tfvars file:
#
# terraform plan -var-file="production.tfvars"

# Example production.tfvars:
# resource_group_name = "prod-rg"
# location            = "eastus"
# server_name         = "prod-sql-server"
# database_name       = "inventory-db"
# database_sku        = "Standard_S2"
# admin_username      = "sqladmin"
# admin_password      = "SecureP@ssw0rd123"  # Use environment variable instead!
# tags = {
#   environment = "production"
#   cost_center = "engineering"
# }
