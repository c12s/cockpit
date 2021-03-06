package helper

const (
	Mutate_NodeConfig = "NodeConfig"
	Mutate_NodeSecret = "NodeSecret"

	ACTIONS    = "Actions"
	SECRETS    = "Secrets"
	CONFIGS    = "Configs"
	NAMESPACES = "Namespaces"
	ROLES      = "Roles"
	USERS      = "users"

	ENV      = "env"
	FILES    = "files"
	AACTIONS = "actions"

	LABELS    = "labels"
	COMPARE   = "compare"
	KIND      = "kind"
	ALL       = "all"
	TYPE      = "type"
	UPDATE    = "update"
	INTERVAL  = "interval"
	FILE      = "file"
	FILE_NAME = "file_name"

	NAMESPACE = "namespace"
	NAME      = "name"
)

// var Allowed_payloads = [...]string{"env", "files", "actions"}
var Configs_payloads = []string{ENV, FILES}
var Actions_payloads = []string{AACTIONS}
var Secrets_payloads = []string{ENV, FILES}
var Namespaces_payloads = []string{LABELS}
