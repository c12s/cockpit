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
