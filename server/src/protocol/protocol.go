package protocol

/*
 *
 */
type GetNodesRequest struct {
	IP string `json:"ip"`
}

type GetNodeMetrixRequest struct {
	IP string `json:"ip"`
	Time string `json:"time"`
}

type GetApplicationInstancesMappingRequest struct {
	Instance string `json:"instance"`
}

type GetNodesMappingRequest struct {
	IP string `json:"ip"`
}

type GetApplicationMetrixRequest struct {
	IP string `json:"ip"`
	Time string `json:"time"`
	Instance string `json:"instance"`
}

/*
 * Page information
 */
type PageInfo struct {
	Begin int `json:"begin" bson:"begin"`
	Number int `json:"number" bson:"number"`
}

func NewPageInfo() *PageInfo {
	return &PageInfo{
		Begin:0,
		Number:0,
	}
}

/*
 * Node using map to save info
 */
type NodeMapInfo struct {
	Info map[string]string `json:"info"`
}

func NewNodeMapInfo() *NodeMapInfo{
	return &NodeMapInfo{
		Info: make(map[string]string),
	}
}

/*
 * Node instance using map to save info
 */
type NodeInstanceMapInfo struct {
	Measurements map[string][]string `json:"measurements"`
}

func NewNodeInstanceMapInfo() *NodeInstanceMapInfo {
	return &NodeInstanceMapInfo{
		Measurements:make(map[string][]string),
	}
}

/*
 * Application instance using map to save info
 */
type ApplicationInstanceMapInfo struct {
	Measurements map[string][]string `json:"measurements"`
}

func NewApplicationInstanceMapInfo() *ApplicationInstanceMapInfo {
	return &ApplicationInstanceMapInfo{
		Measurements:make(map[string][]string),
	}
}

/*
 * Node
 */
type Node struct {
	NodeName string `json:"node_name" bson:"node_name"`
	NodeIP string `json:"node_ip" bson:"node_ip"`

	HostName string `json:"host_name" bson:"host_name"`

	Platform string `json:"platform" bson:"platform"`
	OS string `json:"os" bson:"os"`
	OSVersion string `json:"os_version" bson:"os_version"`
	OSRelease string `json:"os_release" bson:"os_release"`

	MaxCPUs	string `json:"max_cpus" bson:"max_cpus"`
	NCPUs string `json:"ncpus" bson:"ncpus"`

	Bitwith string `json:"bitwidth" bson:"bitwidth"`

	Time string `json:"time" bson:"time"`
}

/*
 * Application instance
 */
type ApplicationInstance struct {
	Name string `json:"name" bson:"name"`
}

/*
 * Application instance & node mapping
 */
type ApplicationInstanceNodeMapping struct {
	Key string `json:"key" bson:"key"`
	Instance string `json:"instance" bson:"instance"`
	NodeIP string `json:"node_ip" bson:"node_ip"`
}

/*
 * View
 */
type Project struct {
	Name string `json:"project" bson:"project"`
}

type Service struct {
	Project string `json:"project" bson:"project"`
	Name string `json:"service" bson:"service"`
}

type Module struct {
	Project string `json:"project" bson:"project"`
	Service string `json:"service" bson:"service"`
	Name string `json:"module" bson:"module"`
	Instances []string `json:"instances" bson:"instances"`
}
