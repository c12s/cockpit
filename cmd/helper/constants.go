package helper

const (
	Mutate_NodeConfig = "NodeConfig"
	Mutate_NodeSecret = "NodeSecret"
)

var allowed_payloads = [...]string{"env", "file", "action"}
