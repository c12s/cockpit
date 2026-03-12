# Cockpit CLI Documentation

Cockpit CLI is a command-line interface (CLI) tool designed for managing various resources within the Cockpit project. This CLI facilitates user registration, node management, labeling, schema validation, configuration management, and more. The goal of this tool is to provide an easy-to-use and maintainable interface for interacting with the Cockpit project's resources.

## Features

- **User Management**: Register and login to the system.
- **Node Management**: List, query, and claim nodes based on specific criteria.
- **Label Management**: Add, update, and delete labels for nodes.
- **Schema Management**: Create, retrieve, validate, and delete schemas.
- **Configuration Management**: Manage both config groups and standalone configs, including placing, listing, and diffing configurations.
- **Starmap Management**: Manage charts on the starmap registry server, including create, retrieve, update, list, delete, and other related operations.


## Table of Contents
- [Getting Started](#getting-started)
- [Command Reference](#command-reference)
  - [User Management](#user-management)
  - [Node Management](#node-management)
  - [Relationship Management](#relationship-management)
  - [Label Management](#label-management)
  - [Schema Management](#schema-management)
  - [Config Group Management](#config-group-management)
  - [Standalone Config Management](#standalone-config-management)
  - [Node Metrics Management](#node-metrics-management)
  - [Starmap Management](#starmap-management)
- [Contributing](#contributing)
- [License](#license)

### Prerequisites
Ensure you have the following installed:
- [Go](https://golang.org/doc/install) (version 1.16 or higher)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [Docker](https://docs.docker.com/get-docker/) (for running the complete environment)

Additionally, you need to clone and set up the tools repository from [https://github.com/c12s/tools](https://github.com/c12s/tools), which includes scripts to pull all other repositories and start Docker containers.

### Installation
1. Clone the tools repository:

    ```sh
    git clone https://github.com/c12s/tools.git
    ```

2. Navigate to the tools directory and follow the setup instructions in the README.md file of the tools repository to pull all necessary repositories and start Docker containers:

    ```sh
    cd tools
    ./install.sh
    ./start.sh
    ```

   or for Windows:

    ```sh
    ./start-windows.sh
    ```

3. Clone the cockpit repository in the same parent directory where the tools repository is located:

    ```sh
    cd ..
    git clone https://github.com/bunjo01/cockpit.git
    ```

4. Navigate to the cockpit project directory:

    ```sh
    cd cockpit
    ```

5. Build the CLI:

    ```sh
    go build -o cockpit
    ```

6. Add the executable to your PATH:
  - On Linux/macOS:

    ```sh
    export PATH=$PATH:$(pwd)
    ```

  - On Windows (Command Prompt):

    ```sh
    set PATH=%PATH%;%cd%
    ```

  - On Windows (PowerShell):

    ```sh
    $env:PATH += ";$PWD"
    ```

## Command Reference

### User Management

#### Register
Register a new user.
- **Command**: cockpit register
- **Options**:
  - --email: Email address of the user.
  - --name: Name of the user.
  - --org: Organization of the user.
  - --surname: Surname of the user.
  - --username: Username of the user.
- **Example**:

    ```sh
    cockpit register --email 'user@gmail.com' --name 'name' --org 'c12s' --surname 'surname' --username 'user'
    ```

#### Login
Login with an existing user.
- **Command**: cockpit login
- **Options**:
  - --username: Username of the user.
- **Example**:

    ```sh
    cockpit login --username 'user'
    ```

### Node Management

#### List Nodes
List all nodes.
- **Command**: cockpit list nodes
- **Options**:
  - --query: Query to filter nodes.
- **Example**:

    ```sh
    cockpit list nodes
    cockpit list nodes --query 'memory-totalGB > 2'
    ```

#### Claim Nodes
Claim nodes based on specific criteria.
- **Command**: cockpit claim nodes
- **Options**:
  - --org: Organization.
  - --query: Query to filter nodes.
- **Example**:

    ```sh
    cockpit claim nodes --org 'c12s' --query 'memory-totalGB > 2'
    ```

#### List Allocated Nodes
List nodes allocated to an organization.
- **Command**: cockpit list nodes allocated
- **Options**:
  - --org: Organization.
  - --query: Query to filter nodes.
- **Example**:

    ```sh
    cockpit list nodes allocated --org 'c12s'
    cockpit list nodes allocated --org 'c12s' --query 'memory-totalGB > 2'
    ```

### Relationship Management

#### Create Relations
Create relations between entities.
- **Command**: cockpit create relations
- **Options**:
  - --ids: IDs of the entities.
  - --kinds: Kinds of the entities.
- **Example**:

    ```sh
    cockpit create relations --ids 'c12s|dev' --kinds 'org|namespace'
    ```

### Label Management

#### Add Label
Add a label to a node.
- **Command**: cockpit put label
- **Options**:
  - --org: Organization of the node.
  - --node-id: ID of the node.
  - --key: Key of the label.
  - --value: Value of the label.
- **Example**:

    ```sh
    cockpit put label --org 'c12s' --node-id 'nodeId' --key 'newlabel' --value '25.00'
    ```

#### Delete Label
Delete a label from a node.
- **Command**: cockpit delete label
- **Options**:
  - --org: Organization of the node.
  - --nodeId: ID of the node.
  - --key: Key of the label.
- **Example**:

    ```sh
    cockpit delete label --org 'c12s' --node-id 'nodeId' --key 'newlabel'
    ```

### Schema Management

#### Create Schema
Create a new schema.
- **Command**: cockpit create schema
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --schema-name: Name of the schema.
  - --version: Version of the schema.
  - --path: Path to the schema file.
- **Example**:

    ```sh
    cockpit create schema --org 'c12s' --namespace 'default' --schema-name 'schema' --version 'v1.0.0' --path 'request/schema/create-schema.yaml'
    ```

#### Get Schema Version
Retrieve schema version details.
- **Command**: cockpit get schema version
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --schema: Name of the schema.
- **Example**:

    ```sh
    cockpit get schema version --org 'c12s' --namespace 'default' --schema-name 'schema'
    ```

#### Get Schema
Retrieve a specific schema.
- **Command**: cockpit get schema
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --schema: Name of the schema.
  - --version: Version of the schema.
- **Example**:

    ```sh
    cockpit get schema --org 'c12s' --namespace 'default' --schema-name 'schema' --version 'v1.0.0'
    ```

#### Validate Schema
Validate a schema.
- **Command**: cockpit validate schema
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --schema: Name of the schema.
  - --version: Version of the schema.
  - --path: Path to the validation file.
- **Example**:

    ```sh
    cockpit validate schema --org 'c12s' --namespace 'default' --schema-name 'schema' --version 'v1.0.0' --path 'request/schema/validate-schema.yaml'
    ```

#### Delete Schema
Delete a schema.
- **Command**: cockpit delete schema
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --schema: Name of the schema.
  - --version: Version of the schema.
- **Example**:

    ```sh
    cockpit delete schema --org 'c12s' --namespace 'default' --schema-name 'schema' --version 'v1.0.0'
    ```

### Config Group Management

#### Add Config Group
Add a configuration group.
- **Command**: cockpit put config group
- **Options**:
  - --path: Path to the config group file.
- **Example**:

    ```sh
    cockpit put config group --path 'request/config-group/create-config-group.yaml'
    ```

#### Get Config Group
Retrieve a configuration group.
- **Command**: cockpit get config group
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --name: Name of the config group.
  - --version: Version of the config group.
  - --output: Output format (json, yaml).
- **Example**:

    ```sh
    cockpit get config group --org 'c12s' --namespace 'default' --name 'app_config' --version 'v1.0.1'
    ```

#### List Config Groups
List all configuration groups.
- **Command**: cockpit list config group
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --output: Output format (json, yaml).
- **Example**:

    ```sh
    cockpit list config group --org 'c12s' --namespace 'default'
    ```

#### Diff Config Groups
Compare differences between configuration groups.
- **Command**: cockpit diff config group
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --names: Names of the config groups.
  - --versions: Versions of the config groups.
  - --output: Output format (json, yaml).
- **Example**:

    ```sh
    cockpit diff config group --org 'c12s' --namespace 'default' --names 'app_config|app_config' --versions 'v1.0.0|v1.0.1'
    cockpit diff config group --org 'c12s' --namespace 'default' --names 'app_config' --versions 'v1.0.0|v1.0.1'
    cockpit diff config group --org 'c12s' --namespace 'default' --names 'app_config|app_config' --versions 'v1.0.0'
    ```

#### Place Config Group
Place a configuration group.
- **Command**: cockpit place config group
- **Options**:
  - --path: Path to the placement file.
- **Example**:

    ```sh
    cockpit place config group --path 'request/config-group/create-config-group-placements.yaml'
    ```

#### List Config Group Placements
List all placements of a configuration group.
- **Command**: cockpit list config group placements
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --name: Name of the config group.
  - --version: Version of the config group.
- **Example**:

    ```sh
    cockpit list config group placements --org 'c12s' --namespace 'default' --name 'app_config' --version 'v1.0.0'
    ```

#### Delete Config Group
Delete a configuration group.
- **Command**: cockpit delete config group
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --name: Name of the config group.
  - --version: Version of the config group.
  - --output: Output format (json, yaml).
- **Example**:

    ```sh
    cockpit delete config group --org 'c12s' --namespace 'default' --name 'app_config' --version 'v1.0.1'
    ```

### Standalone Config Management

#### Add Standalone Config
Add a standalone configuration.
- **Command**: cockpit put standalone config
- **Options**:
  - --path: Path to the config file.
- **Example**:

    ```sh
    cockpit put standalone config --path 'request/standalone-config/create-standalone-config.json'
    ```

#### Get Standalone Config
Retrieve a standalone configuration.
- **Command**: cockpit get standalone config
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --name: Name of the config.
  - --version: Version of the config.
  - --output: Output format (json, yaml).
- **Example**:

    ```sh
    cockpit get standalone config --org 'c12s' --namespace 'default' --name 'db_config' --version 'v1.0.1'
    ```

#### List Standalone Configs
List all standalone configurations.
- **Command**: cockpit list standalone config
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --output: Output format (json, yaml).
- **Example**:

    ```sh
    cockpit list standalone config --org 'c12s' --namespace 'default'
    ```

#### Diff Standalone Configs
Compare differences between standalone configurations.
- **Command**: cockpit diff standalone config
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --names: Names of the configs.
  - --versions: Versions of the configs.
  - --output: Output format (json, yaml).
- **Example**:

    ```sh
    cockpit diff standalone config --org 'c12s' --namespace 'default' --names 'db_config|db_config' --versions 'v1.0.1|v1.0.0'
    cockpit diff standalone config --org 'c12s' --namespace 'default' --names 'db_config' --versions 'v1.0.1|v1.0.0'
    cockpit diff standalone config --org 'c12s' --namespace 'default' --names 'db_config|db_config' --versions 'v1.0.1'
    ```

#### Place Standalone Config
Place a standalone configuration.
- **Command**: cockpit place standalone config
- **Options**:
  - --path: Path to the placement file.
- **Example**:

    ```sh
    cockpit place standalone config --path 'request/standalone-config/create-standalone-config-placements.yaml'
    ```

#### List Standalone Config Placements
List all placements of a standalone configuration.
- **Command**: cockpit list standalone config placements
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --name: Name of the config.
  - --version: Version of the config.
- **Example**:

    ```sh
    cockpit list standalone config placements --org 'c12s' --namespace 'default' --name 'db_config' --version 'v1.0.0'
    ```

#### Delete Standalone Config
Delete a standalone configuration.
- **Command**: cockpit delete standalone config
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --name: Name of the config.
  - --version: Version of the config.
  - --output: Output format (json, yaml).
- **Example**:

    ```sh
    cockpit delete standalone config --org 'c12s' --namespace 'default' --name 'db_config' --version 'v1.0.1'
    ```

### Node Metrics Management

#### Get Node Metrics
Retrieve metrics for a specific node.
- **Command**: cockpit get node metrics
- **Options**:
  - --node-id: Node ID (required).
  - --all: Display all metrics (optional).
  - --sort: Sort metrics by 'cpu', 'memory', 'disk', 'network receive', 'network transmit' or 'bandwidth'.
- **Example**:

    ```sh
    cockpit get node metrics --node-id 'nodeID'
    cockpit get node metrics --node-id 'nodeID' --all-services --sort 'memory'
    ```

### Starmap Management

#### Put Chart
Storing a new chart in the registry.
- **Command**: cockpit starmap put
- **Options**:
  - `--path`: path to file (required).
- **Example**:

    ```sh
    cockpit starmap put --path 'request/starmap/starchart1.yaml'
    ```

#### Get Chart by Metadata
Retrieves a complete chart based on provided metadata.
- **Command**: cockpit starmap metadata
- **Options**:
  - `--maintainer`: chart maintainer (required).
  - `--name`: chart name (required).
  - `--namespace`: chart namespace (optional).
  - `--schema-version`: chart version (optional).
- **Example**:

    ```sh
    cockpit starmap metadata \
    --maintainer 'Nikola' \
    --name 'My chart2' \
    --namespace 'namespace1' \
    —-schema-version 'v1.0.0'

    //Retrieving latest version
    cockpit starmap metadata \
    --maintainer 'Nikola' \
    --name 'My chart2' \
    --namespace 'namespace1'
    ```
#### Get Chart by Id
Retrieves a complete chart based on provided id.
- **Command**: cockpit starmap id
- **Options**:
  - `--maintainer`: chart maintainer (required).
  - `--id`: chart id (required).
  - `--namespace`: chart namespace (optional).
  - `--schema-version`: chart version (optional).
- **Example**:

    ```sh
    cockpit starmap metadata \
    --maintainer 'Nikola' \
    --id 'a0715e59-49c3-4323-bf34-fbb4ccb7e8r4' \
    --namespace 'namespace1' \
    —-schema-version 'v1.0.0'

    # Retrieving latest version
    cockpit starmap metadata \
    --maintainer 'Nikola' \
    --id 'a0715e59-49c3-4323-bf34-fbb4ccb7e8r4' \
    --namespace 'namespace1'
    ```

#### Get Charts by Labels
Retrieves charts based on provided namespace, maintainer, labels, and schema version.

- **Command**: `cockpit starmap labels`
- **Options**:
  - `--maintainer`: Chart maintainer (required)
  - `--namespace`: Chart namespace (optional)
  - `--labels`: Chart labels in `key=value` format (required)
  - `--schema-version`: Chart version (optional)

- **Example**:

```sh
cockpit starmap labels \
  --maintainer 'Nikola' \
  --namespace 'namespace1' \
  --labels key=value \
  --schema-version 'v1.0.0'

  # Retrieving latest version
cockpit starmap labels \
  --maintainer 'Nikola' \
  --namespace 'namespace1' \
  --labels key=value

```
> Multiple labels can be provided by repeating the `--labels` option.

#### Get Missing Layers
Retrieves chart layers that are present in the registry but missing from the provided layer hashes.

- **Command**: `cockpit starmap mslayers`
- **Options**:
  - `--id`: Chart ID (required)
  - `--maintainer`: Chart maintainer (required)
  - `--namespace`: Chart namespace (optional)
  - `--layers`: Layer hash (can be specified multiple times) (required)
  - `--schema-version`: Chart version (optional)

- **Example**:

```sh
cockpit starmap mslayers \
  --id 'a0715e59-49c3-4323-bf34-fbb4ccb7e8r4' \
  --maintainer 'Nikola' \
  --namespace 'namespace1' \
  --layers 53f68b0bdd1e30158bfdbba8925103139181fb0cd6384539dbc8917be15cd969 \
  --layers b2a73f5a23367a37fd6bde5b4aa2d0736ef32f04371f93315949bbe773b7e737 \
  --schema-version 'v1.0.0'
```

#### Get Charts
Retrieves all non-private charts from registry.

- **Command**: `cockpit starmap charts`

- **Example**:

```sh
cockpit starmap charts
```


#### Delete Chart
Deletes a chart from the registry and cleans up orphaned nodes from the graph database.

- **Command**: `cockpit starmap delete`
- **Options**:
  - `--id`: Chart ID (required)
  - `--name`: Chart name (required)
  - `--namespace`: Chart namespace (required)
  - `--maintainer`: Chart maintainer (required)
  - `--kind`: Chart kind (required)
  - `--schema-version`: Chart version (optional)

- **Example**:

```sh
cockpit starmap delete \
  --id 'a0715e59-49c3-4323-bf34-fbb4ccb7e8r4' \
  --maintainer 'Nikola' \
  --kind 'StarChart' \
  --namespace 'namespace1' \
  --name 'My chart4' \
  --schema-version 'v1.0.0'
```

#### Update Chart
Updates a chart’s metadata, version, and associated resources in the registry.

- **Command**: `cockpit starmap update`
- **Options**:
  - `--path`: Path to the chart definition file (JSON or YAML) (required)

- **Example**:

```sh
cockpit starmap update \
  --path 'request/starmap/updateStarchart1.json'
```

#### Switch Checkpoint
Compares two versions of the same chart and determines which components must be started, stopped, or downloaded when switching between versions.

- **Command**: `cockpit starmap swchp`
- **Options**:
  - `--id`: Chart ID (required)
  - `--maintainer`: Chart maintainer (required)
  - `--namespace`: Chart namespace (required)
  - `--old-version`: Source chart version (required)
  - `--new-version`: Target chart version (required)
  - `--layers`: Layer hash (can be specified multiple times)

- **Example**:

```sh
cockpit starmap swchp \
  --id 'a0715e59-49c3-4323-bf34-fbb4ccb7e8r4' \
  --maintainer 'Nikola' \
  --namespace 'namespace1' \
  --old-version 'v1.0.0' \
  --new-version 'v1.0.1' \
  --layers 134626abdf9171255f20b50b56e3a0723c55405f58f518fd1bea3a106e7968e6
```

#### Timeline
Returns a chronological view of all versions of a chart.

- **Command**: `cockpit starmap timeline`
- **Options**:
  - `--id`: Chart ID (required)
  - `--namespace`: Chart namespace (optional)
  - `--maintainer`: Chart maintainer (required)

- **Example**:

```sh
cockpit starmap timeline \
  --id 'a0715e59-49c3-4323-bf34-fbb4ccb7e8r4' \
  --namespace 'namespace1' \
  --maintainer 'Nikola'
```

#### Extend Chart
Extends an existing chart version by creating a new version node and adding only resources that do not already exist in the base version.

- **Command**: `cockpit starmap extend`
- **Options**:
  - `--path`: Path to the chart definition file (JSON or YAML) (required)

- **Example**:

```sh
cockpit starmap extend \
  --path 'request/starmap/extendChart.json'
```

#### Search Components

Searches for components across all charts based on name, description, and tags. Results include matching DataSources, StoredProcedures, EventTriggers, and Events.

- **Command**: `cockpit starmap search`
- **Options**:
  - `--name`: Filter by component name (optional)
  - `--description`: Filter by component description (optional)
  - `--tag`: Filter by tag key=value pair, can be repeated for multiple tags (optional)

- **Example**:
```sh
# Search by name
cockpit starmap search --name 'event2'

# Search by description
cockpit starmap search --description 'some description'

# Search by tags
cockpit starmap search --tag env=prod --tag region=us-east

# Combined search
cockpit starmap search \
  --name 'event2' \
  --description 'some description' \
  --tag env=prod
```

## Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch:

    ```sh
    git checkout -b feature/YourFeature
    ```

3. Make your changes.
4. Commit your changes:

    ```sh
    git commit -m 'Add new feature'
    ```

5. Push to the branch:

    ```sh
    git push origin feature/YourFeature
    ```

6. Open a pull request and describe your changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
