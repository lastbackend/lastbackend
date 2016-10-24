package k8s

// Interface exposes methods on k8s resources.
type LBClientsetInterface interface {
	LB() LBClientInterface
}

type LBClientset struct {
	*LBClient
}

func (c *LBClientset) LB() LBClientInterface {
	if c == nil {
		return nil
	}
	return c.LBClient
}
