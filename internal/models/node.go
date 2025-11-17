package models

// NodeData 表示节点数据
type NodeData struct {
	Name          string `json:"name"`
	Running       string `json:"running"`
	Finished      string `json:"finished"`
	State         string `json:"state"`
	CPUNum        string `json:"cpuNum"`
	MemNum        string `json:"memNum"`
	UpdateTime    string `json:"updateTime"`
	MaxTaskNum    string `json:"maxTaskNum"`
	Version       string `json:"version"`
	ModulesConfig string `json:"modulesConfig"`
}

type Message struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
}
