package minirest

import "github.com/julienschmidt/httprouter"

type handleChain struct {
	handles []handleToHandle
}

type handleToHandle func(httprouter.Handle) httprouter.Handle

func (h *handleChain) handleChain(mainHandle httprouter.Handle) httprouter.Handle {
	for i := len(h.handles) - 1; i > -1; i-- {
		mainHandle = h.handles[i](mainHandle)
	}

	return mainHandle
}
