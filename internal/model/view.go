package model

type View struct {
	Title  string   `json:"title"`
	Fields []string `json:"fields"`
}

// TODO dev
func fillView(m *Model) {
	// if m.Views == nil {
	// 	m.Views = make(map[string]*View)
	// }

	// name := m.Name
	// fields := m.fields
	// if v, ok := m.Views["lists"]; !ok {
	// 	m.Views["lists"] = &View{
	// 		Title:  name + "列表",
	// 		Fields: []string{},
	// 	}
	// } else {
	// 	if v.Title == "" {
	// 		v.Title = name + "列表"
	// 	}
	// 	if len(v.Fields) == 0 {
	// 		v.Fields = fields
	// 	}
	// }
}
