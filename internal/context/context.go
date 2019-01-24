package context

import (
	"reflect"
)

type UserContext struct {
	ContextType       reflect.Type
	IsCookies         bool
	IsHeaders         bool
	RouteParams       map[int]int
	FormParams        map[string]int
	QueryParams       map[string]int
	ExpectJSONRequest bool
}

func NewUserContext() *UserContext {
	return &UserContext{
		IsCookies:         false,
		IsHeaders:         false,
		RouteParams:       make(map[int]int),
		FormParams:        make(map[string]int),
		QueryParams:       make(map[string]int),
		ExpectJSONRequest: false,
	}
}

func (ctx *UserContext) CacheParamsOffset(contextT reflect.Type, routes []string) {
	ctx.ContextType = contextT
	routeParams := make(map[string]int)
	for i := 0; i < contextT.NumField(); i++ {
		tagOfField := contextT.Field(i).Tag
		key, ok := tagOfField.Lookup("route")
		if ok {
			routeParams[key] = i
		}
		key, ok = tagOfField.Lookup("form")
		if ok {
			ctx.FormParams[key] = i
		}
		key, ok = tagOfField.Lookup("query")
		if ok {
			ctx.QueryParams[key] = i
		}
		_, ok = tagOfField.Lookup("json")
		if !ctx.ExpectJSONRequest && ok {
			ctx.ExpectJSONRequest = ok
		}
	}

	for idx, r := range routes {
		// a route part like `:name`
		if r[0] == ':' || r[0] == '*' {
			// r[1:] is `name`, that's the key we expected
			param := r[1:]
			if _, ok := routeParams[param]; ok {
				ctx.RouteParams[idx] = routeParams[param]
			}
		}
	}
}