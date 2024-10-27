package gee

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// 最初handler中的两个参数
	Writer http.ResponseWriter
	Req    *http.Request
	// http请求相关信息
	Path   string
	Method string
	Params map[string]string
	// 响应状态
	StatusCode int
	// 中间件
	handlers []HandlerFunc
	index    int
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	log.Println("[do Next()]")
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// PostForm 查询 Post 请求的参数值
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 查询 GET 请求的参数值
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 设置响应消息的状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 给响头体设置数据
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String 设置纯文本格式的响应消息并写入数据以及响应状态
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 设置 json 格式的响应消息并写入数据以及响应状态
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	// 创建一个 JSON 编码器，并制定将编码结果写入 c.Writer
	encoder := json.NewEncoder(c.Writer)
	// 将 obj 参数编码为 JSON 格式，并写到响应六 c.Writer 中
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data 写入字节数组类型的数据以及响应状态
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML 设置 html 格式的响应消息并写入数据以及响应状态
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
