package model

type TextMessage struct {
	ToUser  string `json:"touser"`
	MsgType string `json:"msgtype"`
	AgentId string `json:"agentid"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

type TextCardMessage struct {
	ToUser   string `json:"touser"`
	MsgType  string `json:"msgtype"`
	AgentId  string `json:"agentid"`
	TextCard struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Url         string `json:"url"`
		Detail      string `json:"btntxt"`
	} `json:"textcard"`
}
