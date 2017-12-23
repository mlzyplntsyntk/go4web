package controller

type BlogController struct {
}

func (t BlogController) SetPattern() string {
	return "blog"
}

func (t BlogController) SetView() string {
	return "blog"
}

func (t BlogController) SetModel() map[string]string {
	return map[string]string{
		"Name": "Blog",
	}
}
