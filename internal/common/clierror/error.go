package clierror

type WithText struct {
	Txt string
}

func (e WithText) Error() string {
	return e.Txt
}
