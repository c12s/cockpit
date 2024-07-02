#!/bin/bash

# Function to execute a command and check its success
run_command() {
    local action=$1
    local entity=$2
    local options=$3
    local fail_msg=$4
    local success_msg=$5

    local cmd="go run main.go $action $entity $options"

    echo "Executing command: $cmd"  # Debugging line
    if eval "$cmd"; then
        echo "$success_msg"
    else
        echo "$fail_msg"
    fi
    echo "==================================================================="
    echo ""
}

# Register with the new user and ensure the context is set
run_command "register" "" "--email 'user@gmail.com' --name 'pera' --org 'c12s' --surname 'surname' --username 'user'" \
            "Register failed" \
            "Register command completed successfully"

# Login with the new user and ensure the context is set
run_command "login" "" "--username user" \
            "Login failed" \
            "Login command completed successfully"

# List nodes
run_command "list" "nodes" "" \
            "List nodes command failed" \
            "List nodes command completed successfully"

run_command "list" "nodes" "--query 'memory-totalGB > 2'" \
            "List nodes with query command failed" \
            "List nodes with query command completed successfully"

# Claim nodes
run_command "claim" "nodes" "--org c12s --query 'memory-totalGB > 2'" \
            "Claim nodes command failed" \
            "Claim nodes command completed successfully"

# List allocated nodes
run_command "list" "nodes allocated" "--org c12s" \
            "List allocated nodes command failed" \
            "List allocated nodes command completed successfully"

run_command "list" "nodes allocated" "--org c12s --query 'memory-totalGB > 2'" \
            "List allocated nodes with query command failed" \
            "List allocated nodes with query command completed successfully"

# List nodes and capture node IDs
node_ids=($(go run main.go list nodes allocated --org c12s | awk '/^[0-9a-f-]+/ {print $1}'))

# Check if at least two nodes are found
if [ ${#node_ids[@]} -lt 2 ]; then
    echo "Not enough nodes found"
    exit 1
fi

# Assign node IDs to variables
node_id_1=${node_ids[0]}
node_id_2=${node_ids[1]}

# Put labels using the captured node IDs
run_command "put" "label" "--org 'c12s' --key 'newlabel' --value '25.00' --node-id $node_id_1" \
            "Put label (newLabel1) command failed" \
            "Put label (newLabel1) command completed successfully"

run_command "put" "label" "--org 'c12s' --key 'newLabel2' --value 'stringValue' --node-id $node_id_1" \
            "Put label (newLabel2) command failed" \
            "Put label (newLabel2) command completed successfully"

run_command "put" "label" "--org 'c12s' --key 'newLabel3' --value 'true' --node-id $node_id_2" \
            "Put label (newLabel3) command failed" \
            "Put label (newLabel3) command completed successfully"

run_command "delete" "label" "--org 'c12s' --key 'newLabel2' --node-id $node_id_1" \
            "Delete label (newLabel2) command failed" \
            "Delete label (newLabel2) command completed successfully"

# CREATE POLICIES
run_command "create" "policies" "--path 'request/policy/create-policy.yaml'" \
            "Create policies command failed" \
            "Create policies command completed successfully"

# CREATE RELATIONS
run_command "create" "relations" "--ids 'c12s|dev' --kinds 'org|namespace'" \
            "Create relations command failed" \
            "Create relations command completed successfully"

# Schema operations
run_command "create" "schema" "--org 'c12s' --schema-name 'schema' --version 'v1.0.0' --path 'request/schema/create-schema.yaml'" \
            "Create schema command failed" \
            "Create schema command completed successfully"

run_command "get" "schema version" "--org c12s -n schema" \
            "Get schema version command failed" \
            "Get schema version command completed successfully"

run_command "get" "schema" "--org c12s -n schema -v v1.0.0" \
            "Get schema command failed" \
            "Get schema command completed successfully"

run_command "delete" "schema" "--org c12s -n schema -v v1.0.0" \
            "Delete schema command failed" \
            "Delete schema command completed successfully"

run_command "create" "schema" "--org 'c12s' --schema-name 'schema' --version 'v1.0.0' --path 'request/schema/create-schema.yaml'" \
            "Create schema command failed" \
            "Create schema command completed successfully"

run_command "validate" "schema" "--org c12s -n schema -v v1.0.0 -p 'request/schema/validate-schema.yaml'" \
            "Validate schema command failed" \
            "Validate schema command completed successfully"

run_command "create" "relations" "--ids 'c12s|dev' --kinds 'org|namespace'" \
            "Create relations command failed" \
            "Create relations command completed successfully"

# CONFIG GROUP

run_command "put" "config group" "--path 'request/config-group/create-config-group.yaml'" \
            "Put config group command failed" \
            "Put config group command completed successfully"

run_command "put" "config group" "--path 'request/config-group/create-config-group.json'" \
            "Put config group command failed" \
            "Put config group command completed successfully"

run_command "get" "config group" "--org 'c12s' --name 'app_config' --version 'v1.0.1'" \
            "Get config group command failed" \
            "Get config group command completed successfully"

run_command "list" "config group" "--org 'c12s'" \
            "List config group command failed" \
            "List config group command completed successfully"

run_command "place" "config group" "--path 'request/config-group/create-config-group-placements.yaml'" \
            "Place config group command failed" \
            "Place config group command completed successfully"

run_command "list" "config group placements" "--org 'c12s' --name 'app_config' --version 'v1.0.0'" \
            "List config group placements command failed" \
            "List config group placements command completed successfully"

run_command "diff" "config group" "--org 'c12s' --names 'app_config|app_config' --versions 'v1.0.0|v1.0.1'" \
            "Diff config group command failed" \
            "Diff config group command completed successfully"

run_command "delete" "config group" "--org 'c12s' --name 'app_config' --version 'v1.0.1'" \
            "Delete config group command failed" \
            "Delete config group command completed successfully"

# STANDALONE CONFIG

run_command "put" "standalone config" "--path 'request/standalone-config/create-standalone-config.json'" \
            "Put standalone config command failed" \
            "Put standalone config command completed successfully"

run_command "put" "standalone config" "--path 'request/standalone-config/create-standalone-config.yaml'" \
            "Put standalone config command failed" \
            "Put standalone config command completed successfully"

run_command "get" "standalone config" "--org 'c12s' --name 'db_config' --version 'v1.0.1'" \
            "Get standalone config command failed" \
            "Get standalone config command completed successfully"

run_command "list" "standalone config" "--org 'c12s'" \
            "List standalone config command failed" \
            "List standalone config command completed successfully"

run_command "diff" "standalone config" "--org 'c12s' --names 'db_config|db_config' --versions 'v1.0.1|v1.0.0'" \
            "Diff standalone config command failed" \
            "Diff standalone config command completed successfully"

run_command "place" "standalone config" "--path 'request/standalone-config/create-standalone-config-placements.yaml'" \
            "Place standalone config command failed" \
            "Place standalone config command completed successfully"

run_command "list" "standalone config placements" "--org 'c12s' --name 'db_config' --version 'v1.0.0'" \
            "List standalone config placements command failed" \
            "List standalone config placements command completed successfully"

run_command "delete" "standalone config" "--org 'c12s' --name 'db_config' --version 'v1.0.1'" \
            "Delete standalone config command failed" \
            "Delete standalone config command completed successfully"

# Node Metrics
run_command "get" "node metrics" "--node-id $node_id_1" \
            "Get node metrics command failed" \
            "Get node metrics command completed successfully"

echo "All commands executed"
