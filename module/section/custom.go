package section

type Custom struct{}

func NewCustom(payload []byte) (*Custom, error) {
	return &Custom{}, nil
}
