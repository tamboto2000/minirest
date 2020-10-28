package minirest

type Service interface {
	Init() error
	Start() error
	Stop() error
	Clean() error
}
