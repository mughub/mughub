package routes

import (
	"net/http"
	"sync"
)

var mapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]interface{})
	},
}

type TmplCtxFunc func(req *http.Request, m map[string]interface{}) error

func GetTmplCtx(req *http.Request, funcs ...TmplCtxFunc) (m map[string]interface{}, err error) {
	m = mapPool.Get().(map[string]interface{})
	defer mapPool.Put(m)

	for _, f := range funcs {
		err = f(req, m)
		if err != nil {
			break
		}
	}

	return
}
