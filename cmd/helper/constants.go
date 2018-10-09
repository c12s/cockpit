package helper

const (
	Mutate_NodeConfig = "NodeConfig"
	Mutate_NodeSecret = "NodeSecret"

	ACTIONS    = "Actions"
	SECRETS    = "Secrets"
	CONFIGS    = "Configs"
	NAMESPACES = "Namespaces"

	ENV   = "env"
	FILES = "files"

	LABELS  = "labels"
	COMPARE = "compare"
	KIND    = "kind"
	ALL     = "all"
	TYPE    = "type"
	UPDATE  = "update"
	FILE    = "file"
)

var Allowed_payloads = [...]string{"env", "files", "actions"}
var Configs_payloads = []string{"env", "files"}
var Actions_payloads = []string{"actions"}
var Secrets_payloads = []string{"keys", "files"}
