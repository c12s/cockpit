package constants

const (
	LongLoginDesc = `Input your username after that you will be prompted to input your password.
Your token will be saved in the token.txt file, which will be sent with all of your request headers.

Example:
- cockpit login --username "username"`

	LongRegisterDesc = `Register a new user by providing an email, name, organization, surname, and username. 
Once these details are entered, you will be prompted to input your password.

Example:
- cockpit register --email "example@gmail.com" --name "name" --org "org" --surname "surname" --username "username"`

	ClaimNodesLongDesc = `Claims nodes for an organization based on a defined query that specifies criteria like labels.
The command allows the organization to take ownership of nodes that match the provided query criteria.
The query can include conditions based on node labels such as memory, CPU, and other attributes. 

Example:
- cockpit claim nodes --org 'org' --query 'labelKey >, =, !=, or < value'
- cockpit claim nodes --org 'org' --query 'memory-totalGB > 2'`

	CreatePoliciesLongDesc = `This command is for creating security policies based on the input file.
Policies are used to define and enforce security rules within the organization. The input file can be in YAML or JSON format, specifying the policy details.

Example:
- cockpit create policies --path 'path to yaml or json file'`

	CreateAppLongDesc = `This command is for creating app specification based on the input file.

Example:
- cockpit create app --path 'path to yaml or json file'`

	CreateNamespaceLongDesc = `This command is for creating namespaces based on the input file.

Example:
- cockpit create namespace --path 'path to yaml or json file'`

	CreateRelationsLongDesc = `This command creates relations between entities specified by their IDs and kinds.
Relations help to establish a hierarchical or dependency structure between different entities within the organization. 
This can include relationships between organizations, namespaces, and other resources.

Example:
- cockpit create relations --ids 'myOrg|dev' --kinds 'org|namespace'`

	CreateSchemaLongDesc = `Creates a schema for an organization by providing schema details and the path to a YAML or JSON file containing the schema definition.
Schemas define the structure of configuration data that can be used across various services and applications within the organization. This command uploads and saves the schema to the server.

Example:
- cockpit create schema --org 'org' --schema-name 'schema' --version 'v1.0.0' --path 'path to yaml or json file'`

	CreateStarmapLongDesc = `This command registers a StarChart definition based on the input file.
The input file must be in YAML/JSON format and contains the complete chart specification.

Example:
- cockpit put starmap --path 'path to starmap.yaml'`

	UpdateStarmapLongDesc = `This command updates an existing StarChart definition in the registry.
The update is performed using an input YAML/JSON file that contains the full StarChart specification.
The chart is identified using its metadata (such as name, namespace, maintainer, and schema version).
If a chart with the specified metadata already exists, it will be updated with the provided definition.

Example:
- Update an existing StarChart:
  cockpit starmap update --path 'path to starmap.yaml'`

	ExtendStarmapLongDesc = `This command extends an existing StarChart version by creating a new version.
Only components and nodes that do not already exist in the base version are added, ensuring that unchanged parts are reused.

Example:
- cockpit starmap extend --path 'path to updated-starmap.yaml'`

	DeleteConfigGroupLongDesc = `This command deletes a specified configuration group version.
The user can specify the organization, the configuration group name, and the version to be deleted.

Example:
- cockpit delete config group --org 'org' --name 'app_config' --version 'v1.0.0'`

	DeleteAppLongDesc = `This command deletes a specified app.
The user must specify the organization, namespace and the application name.

Example:
- cockpit delete app --org 'org' --namespace 'namespace' --name 'app'`

	DeleteNamespaceLongDesc = `This command deletes a specified namespace.
The user must specify the organization and name of the namespace.

Example:
- cockpit delete namespace --org 'org' --name 'name'`

	DeleteNodeLabelsLongDesc = `Delete a specific label from a node using its key.
This command allows the user to remove a label from a node by specifying the node ID, organization, and the label key.

Example:
- cockpit delete label --node-id 'nodeID' --org 'org' --key 'labelKey'`

	DeleteSchemaLongDesc = `This command deletes a schema version from the specified organization.
The user must provide the organization name, schema name, and version to delete the schema. This ensures that the specified schema version is removed from the system.

Example:
- cockpit delete schema --org 'c12s' --schema-name 'schema' --version 'v1.0.1'`

	DeleteStandaloneConfigLongDesc = `This command deletes a specified standalone configuration version.
The user can specify the organization, standalone configuration name, and version to delete the configuration.

Example:
- cockpit delete standalone config --org 'c12s' --name 'db_config' --version 'v1.0.1'`

	DeleteStarmapChartLongDesc = `This command deletes a specific chart from the registry.
The chart can be identified either by its unique chart ID or by providing its name, namespace, maintainer and schema version.

Example:

- Delete a specific chart version:
  cockpit starmap delete --id 'a0715e59-49c3-4323-bf34-fbb4ccb7e8r4' --name 'myChart' --namespace 'dev' --maintainer 'org' --schema-version 'v1.0.1' --kind 'StarMap'`

	DiffConfigGroupLongDesc = `This command compares two configuration groups specified by their names and versions, displays the differences, and saves them to both YAML and JSON files (optional).
The user can specify the organization, names, and versions of the two configuration groups to be compared.

Example:
- cockpit diff config group --org 'org' --names 'name1|name2' --versions 'version1|version2'
- cockpit diff config group --org 'org' --names 'name' --versions 'version1|version2'
- cockpit diff config group --org 'org' --names 'name1|name2' --versions 'version'`

	DiffStandaloneConfigLongDesc = `This command compares two standalone configurations specified by their names and versions, displays the differences, and saves them to both YAML and JSON files (optional).
The user can specify the organization, names, and versions of the two configuration groups to be compared.

Example:
- cockpit diff standalone config --org 'org' --names 'name1|name2' --versions 'version1|version2'
- cockpit diff standalone config--org 'org' --names 'name' --versions 'version1|version2'
- cockpit diff standalone config --org 'org' --names 'name1|name2' --versions 'version'`

	StarmapSwitchCheckpointLongDesc = `This command compares two versions of the same StarChart and determines the required changes when switching between them.
It analyzes which components must be started, stopped, or downloaded in order to move from the old chart version to the new one.
Optionally, a list of already available layer hashes can be provided to avoid unnecessary downloads.

Example:
- Compare two chart versions and determine required actions:
  cockpit starmap swchp \
    --id 'a0715e59-49c3-4323-bf34-fbb4ccb7e8r4' \
    --maintainer 'Nikola' \
    --namespace 'namespace1' \
    --old-version 'v1.0.0' \
    --new-version 'v1.0.1' \
    --layer 53f68b0bdd1e30158bfdbba8925103139181fb0cd6384539dbc8917be15cd969`

	GetAppConfigLongDesc = `This command retrieves a specific configuration by its organization, name, and version.
The user can specify the organization, configuration name, and version to retrieve the configuration details. The response can be formatted as either YAML or JSON based on user preference.

Example:
- cockpit get config group --org 'org' --name 'app_config' --version 'v1.0.0'`

	GetStandaloneConfigLongDesc = `This command retrieves a specific standalone configuration by its organization, name, and version.
The user can specify the organization, configuration name, and version to retrieve the configuration details. The response can be formatted as either YAML or JSON based on user preference.

Example:
- cockpit get standalone configuration --org 'org' --name 'db_config' --version 'v1.0.0'`

	LatestMetricsLongDesc = `This command fetches the latest metrics for a specific node and displays them.
The user can specify the node ID to retrieve the metrics. 

Example:
- cockpit get nodes metrics --node-id 'nodeID'`

	GetSchemaLongDesc = `This command retrieves the schema from a specified organization and specific version and saves it to a YAML or JSON file (optional).
The user can specify the organization, schema name, and version to retrieve the schema details.

Example:
- cockpit get schema --org 'org' --schema-name 'schema_name' --version 'v1.0.0'`

	GetSchemaVersionLongDesc = `This command retrieves schema versions for a specific schema.
The user can specify the organization and schema name to retrieve the list of schema versions.

Example:
- cockpit get schema version --org 'org' --schema-name 'schema_name'`

	GetStarmapMetadataLongDesc = `This command retrieves detailed metadata for a specific StarMap chart from the server.
You can specify the chart by providing its name, namespace, maintainer, and optionally the schema version.
If the schema version is omitted, the latest version of the chart will be retrieved automatically.

Example:
- Retrieve a specific version:
  cockpit starmap metadata --name 'myChart' --namespace 'dev' --maintainer 'org' --schema-version 'v1.0.0'
- Retrieve the latest version:
  cockpit starmap metadata --name 'myChart' --namespace 'dev' --maintainer 'org'`

	GetStarmapIdLongDesc = `This command retrieves detailed metadata for a specific StarMap chart from the server.
You can specify the chart by providing its unique chart ID.
Optionally, you may also specify the schema version to retrieve a particular version of the chart.
If the schema version is omitted, the latest version of the chart will be retrieved automatically.

Example:
- Retrieve a specific version by chart ID:
  cockpit starmap metadata --id 'a0715e59-49c3-4323-bf34-fbb4ccb7e8r4' --maintainer 'org' --namespace 'dev' --schema-version 'v1.0.0'
- Retrieve the latest version by chart ID:
  cockpit starmap metadata --id 'a0715e59-49c3-4323-bf34-fbb4ccb7e8r4' --maintainer 'org' --namespace 'dev'`

	GetStarmapAllChartsLongDesc = `This command retrieves all non-private charts from the StarMap registry.
Example:
- Retrieve all public charts:
  cockpit starmap charts`

	GetStarmapLabelsLongDesc = `This option allows you to filter StarMap charts using labels.
Labels are specified as key-value pairs and can be provided multiple times.
Only charts that match all specified labels will be considered.

Example:
- Retrieve a charts filtered by labels:
  cockpit starmap labels --maintainer 'Nikola' --namespace 'namespace1' --labels key=value`

	GetMissingLayersLongDesc = `This command retrieves chart layers that are present in the registry but missing from the provided layer hashes.
You can specify the chart by providing its unique chart ID, maintainer, namespace, and optional schema version.
You must also provide the list of layer hashes to check against the registry.
Each --layer flag represents one layer hash. You can repeat the flag multiple times to provide multiple hashes.

Example:
- Retrieve missing layers for a specific chart:
  cockpit starmap missing-layers \
    --id 'a0715e59-49c3-4323-bf34-fbb4ccb7e8r4' \
    --maintainer 'Nikola' \
    --namespace 'namespace1' \
    --schema-version 'v1.0.1' \
    --layer 5db611d4ef88f8ff8b8a049a007ded454d59e9a39a68b5d6b72706c068aa45ca \
    --layer e56b6ad61933b816e80366c321271cd36c971521ba72c53dbec20469ba00531f \
    --layer ab1aed768467b8405cee37bf008b4243322370d3d78ef687dcf25920160b1d80`

	GetStarmapChartTimelineLongDesc = `This command returns a chronological timeline of all versions of a StarChart.
It provides an ordered view of chart versions, allowing users to track evolution, changes, and version history over time.

Example:
- cockpit starmap timeline --id 'a0715e59-49c3-4323-bf34-fbb4ccb7e8r4' --namespace 'namespace1' --maintainer 'Nikola'`

	SearchStarmapLongDesc = `This command searches for components across all charts based on name, description, and tags.
You can filter results by providing any combination of name, description, and tag key-value pairs.
Results include matching DataSources, StoredProcedures, EventTriggers, and Events.
Each --tag flag represents one key=value pair. You can repeat the flag multiple times to filter by multiple tags.

Example:
- Search for components by name:
  cockpit starmap search --name 'event2'

- Search by description:
  cockpit starmap search --description 'some description'

- Search by tags:
  cockpit starmap search --tag env=prod --tag region=us-east

- Combined search:
  cockpit starmap search --name 'event2' --description 'some description' --tag env=prod`

	AllocatedNodesLongDesc = `This command allows you to list all nodes allocated to a specified organization.
You can also use a query to search for nodes based on their labels.
The query format allows you to filter nodes using operators like >, =, !=, and < with the label values.

Examples:
- cockpit list nodes allocated --org 'org' --query 'labelKey >, =, !=, or < value'
- cockpit list nodes allocated --org 'org' --query 'memory-totalGB > 2'`

	ListConfigGroupPlacementsLongDesc = `This command retrieves all configuration group placements from a specified organization,
displays them in a nicely formatted way, and allows you to see the placements in detail.

Examples:
- cockpit list config group placements --org 'org' --name 'app_config' --version 'v1.0.0'
- cockpit list config group placements --org 'org' --name 'db_config' --version 'v2.0.0'`

	ListNodesLongDesc = `Retrieve a comprehensive list of all available nodes in the system.
These nodes can be allocated to your organization based on your requirements.
You can use a query to filter the nodes using operators like >, =, !=, and < with the label values.

Examples:
- cockpit list nodes --query 'labelKey >, =, !=, or < value'
- cockpit list nodes --query 'memory-totalGB > 2'`

	ListStandaloneConfigLongDesc = `This command retrieves a list of all standalone configurations for a given organization.

Examples:
- cockpit list standalone config --org 'org' --output 'json'
- cockpit list standalone config --org 'org' --output 'yaml'`

	ListStandaloneConfigPlacementsLongDesc = `This command retrieves all standalone configuration placements from a specified organization,
displays them in a nicely formatted way, and allows you to see the placements in detail.

Examples:
- cockpit list standalone config placements --org 'org' --name 'app_config' --version 'v1.0.0'
- cockpit list standalone config placements --org 'org' --name 'db_config' --version 'v2.0.0'`

	PlaceConfigGroupPlacementsLongDesc = `This command places configuration group placements based on the input file.
The input file should be in either YAML or JSON format, containing the details of the configuration group placements.
It reads the file, processes the placements, and applies them accordingly.

Example:
- cockpit place config group placements --path 'path to yaml or json file'`

	PlaceStandaloneConfigPlacementsLongDesc = `This command places standalone configuration placements based on the input file.
The input file should be in either YAML or JSON format, containing the details of the standalone configuration placements.
It reads the file, processes the placements, and applies them accordingly.

Example:
- cockpit place standalone config placements --path 'path to yaml or json file'`

	PutConfigGroupLongDesc = `This command sends a configuration group read from a file (JSON or YAML) to the server.
It processes the file and uploads the configuration group, displaying the server's response in the same format as the input file.

Example:
- cockpit put config group --path 'path to yaml or JSON file'`

	LongLabelDesc = `This command allows you to add a new label to a specified node, enhancing node metadata.
Provide a key-value pair to define the label. If the label already exists, its value will be updated to the new specified value.
The command supports different types of values: strings, boolean, and floating-point numbers.
The input format determines the appropriate type and URL for the request.

Examples:
- cockpit put label --key 'env' --value 'production' --node-id 'nodeId' --org 'org'
- cockpit put label --key 'active' --value 'true' --node-id 'nodeId' --org 'org'
- cockpit put label --key 'cpu' --value '2.5' --node-id 'nodeId' --org 'org'`

	PutStandaloneConfigLongDesc = `This command sends a standalone configuration read from a file (JSON or YAML) to the server.
It processes the file and uploads the standalone configuration, displaying the server's response in the same format as the input file.

Example:
- cockpit put standalone config --path 'path to yaml or JSON file'`

	ValidateSchemaVersionLongDesc = `This command validates a schema version with the given configuration.
The user specifies the organization, schema name, version, and path to the YAML or JSON configuration file.
It reads the configuration file and validates the schema version against it.

Example:
- cockpit validate schema --org 'org' --schema-name 'schema' --version 'v1.0.0' --path '/path/to/config.yaml'`
)
