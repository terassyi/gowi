package decoder

type custom struct{}

func newCustom(payload []byte) (*custom, error) {
	return &custom{}, nil
}

func (*custom) detail() string {
	return "not implemented."
}
