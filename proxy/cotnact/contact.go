package cotnact

type ContactProxy struct {
	address string
}

func New(address string) *ContactProxy {
	return &ContactProxy{
		address: address,
	}
}
