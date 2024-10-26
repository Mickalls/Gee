package main

import (
	"Gee/gee"
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *gee.Router {
	r := gee.NewRouter()
	r.AddRoute("GET", "/", nil)
	r.AddRoute("GET", "/hello/:name", nil)
	r.AddRoute("GET", "/hello/b/c", nil)
	r.AddRoute("GET", "/hi/:name", nil)
	r.AddRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(gee.ParsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(gee.ParsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(gee.ParsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.GetRoute("GET", "/hello/wangyouhua")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.Pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "wangyouhua" {
		t.Fatal("name should match wangyouhua")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.Pattern, ps["name"])
}
