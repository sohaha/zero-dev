package api

type Route struct {
	Path    string `json:"path"`
	Handler string `json:"handler"`
	Method  string `json:"method"`
}

type HTTP struct {
	Name   string   `json:"name"`
	Routes []Router `json:"routes"`
}

func InitHTTP() []HTTP {
	return []HTTP{
		HTTP{},
	}
}
