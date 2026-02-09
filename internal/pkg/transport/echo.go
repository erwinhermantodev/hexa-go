package transport

import (
	"github.com/erwinhermantodev/hexa-go/internal/pkg/interfaces"
	"github.com/labstack/echo/v4"
)

type echoRouter struct {
	echo *echo.Echo
}

type echoContext struct {
	ctx echo.Context
}

func NewEchoRouter() interfaces.Router {
	return &echoRouter{echo: echo.New()}
}

func (r *echoRouter) Group(prefix string, middleware ...interface{}) interfaces.Router {
	// Middleware conversion logic would go here
	return &echoRouterGroup{group: r.echo.Group(prefix)}
}

func (r *echoRouter) GET(path string, handler interface{}, middleware ...interface{}) {
	h := handler.(func(interfaces.Context) error)
	r.echo.GET(path, func(c echo.Context) error {
		return h(&echoContext{ctx: c})
	})
}

func (r *echoRouter) POST(path string, handler interface{}, middleware ...interface{}) {
	h := handler.(func(interfaces.Context) error)
	r.echo.POST(path, func(c echo.Context) error {
		return h(&echoContext{ctx: c})
	})
}

func (r *echoRouter) PUT(path string, handler interface{}, middleware ...interface{}) {
	h := handler.(func(interfaces.Context) error)
	r.echo.PUT(path, func(c echo.Context) error {
		return h(&echoContext{ctx: c})
	})
}

func (r *echoRouter) DELETE(path string, handler interface{}, middleware ...interface{}) {
	h := handler.(func(interfaces.Context) error)
	r.echo.DELETE(path, func(c echo.Context) error {
		return h(&echoContext{ctx: c})
	})
}

func (r *echoRouter) Serve(address string) error {
	return r.echo.Start(address)
}

// Helpers for Context
func (c *echoContext) Bind(i interface{}) error {
	return c.ctx.Bind(i)
}

func (c *echoContext) JSON(code int, i interface{}) error {
	return c.ctx.JSON(code, i)
}

func (c *echoContext) Param(name string) string {
	return c.ctx.Param(name)
}

func (c *echoContext) QueryParam(name string) string {
	return c.ctx.QueryParam(name)
}

func (c *echoContext) Get(key string) interface{} {
	return c.ctx.Get(key)
}

func (c *echoContext) Set(key string, val interface{}) {
	c.ctx.Set(key, val)
}

// Group helper
type echoRouterGroup struct {
	group *echo.Group
}

func (g *echoRouterGroup) Group(prefix string, middleware ...interface{}) interfaces.Router {
	return &echoRouterGroup{group: g.group.Group(prefix)}
}

func (g *echoRouterGroup) GET(path string, handler interface{}, middleware ...interface{}) {
	h := handler.(func(interfaces.Context) error)
	g.group.GET(path, func(c echo.Context) error {
		return h(&echoContext{ctx: c})
	})
}

func (g *echoRouterGroup) POST(path string, handler interface{}, middleware ...interface{}) {
	h := handler.(func(interfaces.Context) error)
	g.group.POST(path, func(c echo.Context) error {
		return h(&echoContext{ctx: c})
	})
}

func (g *echoRouterGroup) PUT(path string, handler interface{}, middleware ...interface{}) {
	h := handler.(func(interfaces.Context) error)
	g.group.PUT(path, func(c echo.Context) error {
		return h(&echoContext{ctx: c})
	})
}

func (g *echoRouterGroup) DELETE(path string, handler interface{}, middleware ...interface{}) {
	h := handler.(func(interfaces.Context) error)
	g.group.DELETE(path, func(c echo.Context) error {
		return h(&echoContext{ctx: c})
	})
}

func (g *echoRouterGroup) Serve(address string) error {
	return nil // Not applicable for groups
}
