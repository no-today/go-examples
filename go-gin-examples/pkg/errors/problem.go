package errors

import "fmt"

type Problem struct {
	error
	// Faq url, answer why this exception occurs
	Guide string `json:"guide,omitempty"`
	// Requested path
	Path string `json:"path,omitempty"`
	// Response status
	Status int `json:"status,omitempty"`
	// Error title
	Title string `json:"title,omitempty"`
	// Error detail
	Description string `json:"description,omitempty"`
	// Server instance
	Instance string `json:"instance,omitempty"`
	// 4xx Appear, the description request parameter is incorrect
	Parameters interface{} `json:"parameters,omitempty"`
	// Call Stack
	Stack []interface{} `json:"stack,omitempty"`
}

func (p Problem) Error() string {
	return p.Title
}

func (p Problem) Descf(description string, a ...interface{}) Problem {
	p.Description = fmt.Sprintf(description, a...)
	return p
}

func (p Problem) Err(err error) Problem {
	p.error = err
	return p
}

type problemBuild struct {
	guide       string
	path        string
	status      int
	title       string
	description string
	instance    string
	parameters  interface{}
	stack       []interface{}
}

func Build() *problemBuild {
	return &problemBuild{}
}

func (p *problemBuild) Url(url string) *problemBuild {
	p.path = url
	return p
}

func (p *problemBuild) Status(status int) *problemBuild {
	p.status = status
	return p
}

func (p *problemBuild) Title(title string) *problemBuild {
	p.title = title
	return p
}

func (p *problemBuild) Description(description string) *problemBuild {
	p.description = description
	return p
}

func (p *problemBuild) Instance(instance string) *problemBuild {
	p.instance = instance
	return p
}

func (p *problemBuild) Parameters(parameters interface{}) *problemBuild {
	p.parameters = parameters
	return p
}

func (p *problemBuild) Stack(stack ...interface{}) *problemBuild {
	p.stack = stack
	return p
}

func (p *problemBuild) Build() *Problem {
	return &Problem{
		Guide:       p.guide,
		Path:        p.path,
		Status:      p.status,
		Title:       p.title,
		Description: p.description,
		Instance:    p.instance,
		Parameters:  p.parameters,
		Stack:       p.stack,
	}
}
