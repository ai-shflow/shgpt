package hellogpt

import (
	"sync"
)

var (
	Options = []string{
		"option",
	}
)

type Context struct {
	mutex sync.Mutex
	Debug bool
}

func NewContext(debug bool) (*Context, error) {
	return &Context{
		Debug: debug,
	}, nil
}

func (c *Context) Init() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// TBD: FIXME

	return nil
}

func (c *Context) Deinit() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// TBD: FIXME

	return nil
}

func (c *Context) Run(args map[string]string) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// TBD: FIXME
	ret := "hello gpt"

	return ret, nil
}
