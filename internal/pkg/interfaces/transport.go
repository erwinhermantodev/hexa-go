package interfaces

type Router interface {
	Group(prefix string, middleware ...interface{}) Router
	GET(path string, handler interface{}, middleware ...interface{})
	POST(path string, handler interface{}, middleware ...interface{})
	PUT(path string, handler interface{}, middleware ...interface{})
	DELETE(path string, handler interface{}, middleware ...interface{})
	Serve(address string) error
}

type Context interface {
	Bind(i interface{}) error
	JSON(code int, i interface{}) error
	Param(name string) string
	QueryParam(name string) string
	Get(key string) interface{}
	Set(key string, val interface{})
}
