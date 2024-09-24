package configuration

type InputSlice struct {
	Env Env
}

func (islice InputSlice) String() string {
	return islice.Env.String()
}