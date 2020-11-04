package minirest

import "github.com/julienschmidt/httprouter"

type handleChain struct {
	handles []handleToHandle
}

type handleToHandle func(httprouter.Handle) httprouter.Handle

func (h *handleChain) handleChain(mainHandle httprouter.Handle) httprouter.Handle {
	for _, handle := range h.handles {
		mainHandle = handle(mainHandle)
	}

	return mainHandle
}
