package gee

import (
	"net/http"
	"strings"
)

type Router struct {
	handlers map[string]HandlerFunc
	roots    map[string]*trieNode // roots 存储每种请求方式的Trie树根节点
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]HandlerFunc),
		roots:    make(map[string]*trieNode),
	}
}

func (r *Router) AddRoute(method string, pattern string, handler HandlerFunc) {
	parts := ParsePattern(pattern)
	key := method + "-" + pattern
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &trieNode{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *Router) GetRoute(method string, path string) (*trieNode, map[string]string) {
	searchParts := ParsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n == nil {
		return nil, nil
	}
	parts := ParsePattern(n.Pattern)
	for idx, part := range parts {
		if part[0] == ':' {
			params[part[1:]] = searchParts[idx]
		} else if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(searchParts[idx:], "/")
			break
		}
	}
	return n, params
}

func (r *Router) handle(c *Context) {
	n, params := r.GetRoute(c.Method, c.Path)
	if n != nil {
		key := c.Method + "-" + n.Pattern
		c.Params = params
		// 将原本路由对应的handler放在中间件的最后面
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 Not Found: %s\n", c.Path)
		})
	}
	c.Next()
}

func ParsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item == "" {
			continue
		}
		parts = append(parts, item)
		if item[0] == '*' {
			break
		}
	}
	return parts
}
