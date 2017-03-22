package context

var _ctx ctx

func Get() *ctx {
	return &_ctx
}

type ctx struct {

}

