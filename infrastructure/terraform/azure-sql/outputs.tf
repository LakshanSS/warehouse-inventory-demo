# Outputs file for Azure SQL Database
# These outputs are captured by the workflow and saved to Kubernetes secrets

# Output definitions are in main.tf
# This file documents how outputs are used

# After terraform apply, outputs can be accessed:
#
#   terraform output server_fqdn
#   terraform output -raw connection_string
#   terraform output -json > outputs.json
#
# The workflow captures these outputs and creates a Kubernetes secret:
#
#   kubectl get secret <database-name>-connection -o yaml
#
# Secret structure:
#   data:
#     server: <base64-encoded server FQDN>
#     database: <base64-encoded database name>
#     port: <base64-encoded "1433">
#     connection-string: <base64-encoded full connection string>
#
# Usage in workloads:
#
#   env:
#     - name: DB_SERVER
#       valueFrom:
#         secretKeyRef:
#           name: pge-demo-db-connection
#           key: server
