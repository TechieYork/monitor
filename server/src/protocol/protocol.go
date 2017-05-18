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
 * Node
 */
type Node struct {
	Info map[string]string `json:"info"`
}

func NewNode() *Node{
	return &Node{
		Info: make(map[string]string),
	}
}

/*
 * Node in mongo
 */
type NodeInMongo struct {
	Info NodeInfo `json:"info" bson:"info"`
	ApplicationInstances map[string]string `json:"application_instances" bson:"application_instances"`
}

type NodeInfo struct {
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
 * Node instance
 */
type NodeInstance struct {
	Measurements map[string][]string `json:"measurements"`
}

func NewNodeInstance() *NodeInstance {
	return &NodeInstance{
		Measurements:make(map[string][]string),
	}
}

/*
 * Application instance
 */
type ApplicationInstance struct {
	Measurements map[string][]string `json:"measurements"`
}

func NewApplicationInstance() *ApplicationInstance {
	return &ApplicationInstance{
		Measurements:make(map[string][]string),
	}
}

/*
 * Application instance in mongodb
 */
type ApplicationInstanceInMongo struct {
	Info ApplicationInstanceInfo `json:"info" bson:"info"`
}

type ApplicationInstanceInfo struct {
	Name string `json:"name" bson:"name"`
}

/*
 * Application instance & node mapping
 */
type ApplicationInstanceNodeMapping struct {
	Info ApplicationInstanceNodeMappingInfo `json:"info" bson:"info"`
}

type ApplicationInstanceNodeMappingInfo struct {
	Key string `json:"key" bson:"key"`
	Instance string `json:"instance" bson:"instance"`
	NodeIP string `json:"node_ip" bson:"node_ip"`
}
