package web

import (
	"net/http"

	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/service/stack"
)

func NewService(opts ...service.Option) service.Service {
	options := newOptions(opts...)

	options = append(options,
		service.BeforeInit(setHandle),
	)

	return stack.NewService(options...)
}

func setHandle(sOpts *service.Options) error {
	var mux http.Handler
	if sOpts.Context.Value(serverMuxKey{}) != nil {
		if muxTmp, ok := sOpts.Context.Value(serverMuxKey{}).(http.Handler); ok {
			mux = muxTmp
		}
	} else {
		muxTmp := http.NewServeMux()
		// handler funcs
		if sOpts.Context.Value(handlerFuncsKey{}) != nil {
			if handlers, ok := sOpts.Context.Value(handlerFuncsKey{}).([]HandlerFunc); ok {
				for _, handler := range handlers {
					muxTmp.HandleFunc(handler.Route, handler.Func)
				}
			}
		}

		// static dir
		if sOpts.Context.Value(staticDirKey{}) != nil {
			if sd, ok := sOpts.Context.Value(staticDirKey{}).(staticDir); ok {
				muxTmp.Handle(sd.Route, http.StripPrefix(sd.Route, http.FileServer(http.Dir(sd.Dir))))
			}
		}

		mux = muxTmp
	}

	return sOpts.Server.Handle(sOpts.Server.NewHandler(mux))
}
