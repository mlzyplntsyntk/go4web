package controller

type IndexController struct {
}

func (t IndexController) SetPattern() string {
	return "home"
}

func (t IndexController) SetView() string {
	return "home"
}

func (t IndexController) SetModel() map[string]string {
	return map[string]string{
		"Name": "Testing",
	}
}
