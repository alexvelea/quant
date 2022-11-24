package utils

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Float64P(f float64) *float64 {
	return &f
}
