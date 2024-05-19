package clients

const (
	Gateway                     = "http://localhost:5555"
	Api                         = Gateway + "/apis/core/v1"
	LoginEndpoint               = Api + "/auth"
	RegisterEndpoint            = Api + "/users"
	AvailableNodesEndpoint      = Api + "/nodes/available"
	AvailableNodesQueryEndpoint = Api + "/nodes/available/query_match"
	AllocatedNodesEndpoint      = Api + "/nodes/allocated"
	AllocatedNodesQueryEndpoint = AllocatedNodesEndpoint + "/query_match"
	ClaimNodesEndpoint          = Api + "/nodes"
	LabelsEndpoint              = Api + "/labels"
	LabelsFloatEndpoint         = Api + "/labels/float64"
	LabelsStringEndpoint        = Api + "/labels/string"
	LabelsBoolEndpoint          = Api + "/labels/bool"
)
