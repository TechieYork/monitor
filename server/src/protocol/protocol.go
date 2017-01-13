package protocol

/*
 *
 */
type GetNodesRequest struct {
	IP string `json:"ip"`
}

type GetNodeInstancesRequest struct {
	IP string `json:"ip"`
}

type GetNodeMetrixRequest struct {
	IP string `json:"ip"`
	Time string `json:"time"`
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
 * Metrix
 */
type Metrix struct {
	Serie string `json:"serie"`
	Instance string `json:"instance"`
	TimeBegin string `json:"time_begin"`
	Interval string `json:"interval"`
	Points []float64 `json:"points"`
}

func NewMetrix() *Metrix {
	return &Metrix{
		Serie: "null",
		Instance: "null",
		TimeBegin: "null",
		Interval: "null",
		Points: []float64{},
	}
}