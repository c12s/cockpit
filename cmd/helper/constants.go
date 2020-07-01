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
	TOPOLOGY   = "Topology"

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

	B  = "b"
	KB = "kb"
	MB = "mb"
	GB = "gb"
	TB = "tb"
	BV = 1
)

// var Allowed_payloads = [...]string{"env", "files", "actions"}
var Configs_payloads = []string{ENV, FILES}
var Actions_payloads = []string{AACTIONS}
var Secrets_payloads = []string{ENV, FILES}
var Namespaces_payloads = []string{LABELS}

var maper = map[string]int64{
	B:  BV,
	KB: BV << 10,
	MB: BV << 20,
	GB: BV << 30,
	TB: BV << 40,
}
