package protocol

import "time"

type Data struct {
	Time string `json:"time"`
	Tag map[string]interface{} `json:"tag"`
	Field map[string]interface{} `json:"field"`
}

func NewData() *Data{
	return &Data{
		Time: time.Now().String(),
		Tag: make(map[string]interface{}),
		Field: make(map[string]interface{}),
	}
}

type Proto struct {
	Name string `json:"name"`
	Version int `json:"version"`
	DataList []Data `json:"data"`
}

func NewProto(version int) *Proto {
	return &Proto{
		Name: "",
		Version: version,
		DataList: []Data{},
	}
}
