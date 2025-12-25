package router

import (
	"github.com/gin-gonic/gin"
)

// Router represents the main router
type Router struct {
	engine *gin.Engine
	groups []*Group
}

// New creates a new router instance
func New() *Router {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	return &Router{
		engine: engine,
		groups: make([]*Group, 0),
	}
}

// Engine returns the underlying gin engine
func (r *Router) Engine() *gin.Engine {
	return r.engine
}

// Group creates a new route group with prefix
func (r *Router) Group(prefix string) *Group {
	group := &Group{
		RouterGroup: r.engine.Group(prefix),
		prefix:      prefix,
	}
	r.groups = append(r.groups, group)
	return group
}

// GET registers a GET route
func (r *Router) GET(path string, handler gin.HandlerFunc) {
	r.engine.GET(path, handler)
}

// POST registers a POST route
func (r *Router) POST(path string, handler gin.HandlerFunc) {
	r.engine.POST(path, handler)
}

// PUT registers a PUT route
func (r *Router) PUT(path string, handler gin.HandlerFunc) {
	r.engine.PUT(path, handler)
}

// DELETE registers a DELETE route
func (r *Router) DELETE(path string, handler gin.HandlerFunc) {
	r.engine.DELETE(path, handler)
}

// PATCH registers a PATCH route
func (r *Router) PATCH(path string, handler gin.HandlerFunc) {
	r.engine.PATCH(path, handler)
}

// Use adds middleware to the router
func (r *Router) Use(middlewares ...gin.HandlerFunc) {
	r.engine.Use(middlewares...)
}

// Run starts the HTTP server
func (r *Router) Run(addr ...string) error {
	return r.engine.Run(addr...)
}

// Group represents a route group
type Group struct {
	*gin.RouterGroup
	prefix string
}

// Subgroup creates a subgroup within a group
func (g *Group) Subgroup(prefix string) *Group {
	group := &Group{
		RouterGroup: g.RouterGroup.Group(prefix),
		prefix:      g.prefix + prefix,
	}
	return group
}

// GET registers a GET route in the group
func (g *Group) GET(path string, handler gin.HandlerFunc) {
	g.RouterGroup.GET(path, handler)
}

// POST registers a POST route in the group
func (g *Group) POST(path string, handler gin.HandlerFunc) {
	g.RouterGroup.POST(path, handler)
}

// PUT registers a PUT route in the group
func (g *Group) PUT(path string, handler gin.HandlerFunc) {
	g.RouterGroup.PUT(path, handler)
}

// DELETE registers a DELETE route in the group
func (g *Group) DELETE(path string, handler gin.HandlerFunc) {
	g.RouterGroup.DELETE(path, handler)
}

// PATCH registers a PATCH route in the group
func (g *Group) PATCH(path string, handler gin.HandlerFunc) {
	g.RouterGroup.PATCH(path, handler)
}

// Use adds middleware to the group
func (g *Group) Use(middlewares ...gin.HandlerFunc) {
	g.RouterGroup.Use(middlewares...)
}
