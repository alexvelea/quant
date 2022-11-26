package utils

func PanicIf(value bool, err error) {
	if value == true {
		panic(err)
	}
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Float64P(f float64) *float64 {
	return &f
}
