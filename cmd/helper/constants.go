package helper

const (
	Mutate_NodeConfig = "NodeConfig"
	Mutate_NodeSecret = "NodeSecret"

	ACTIONS    = "Actions"
	SECRETS    = "Secrets"
	CONFIGS    = "Configs"
	NAMESPACES = "Namespaces"

	ENV      = "env"
	FILES    = "files"
	AACTIONS = "actions"

	LABELS  = "labels"
	COMPARE = "compare"
	KIND    = "kind"
	ALL     = "all"
	TYPE    = "type"
	UPDATE  = "update"
	FILE    = "file"
)

// var Allowed_payloads = [...]string{"env", "files", "actions"}
var Configs_payloads = []string{ENV, FILES}
var Actions_payloads = []string{AACTIONS}
var Secrets_payloads = []string{ENV, FILES}
var Namespaces_payloads = []string{LABELS}
