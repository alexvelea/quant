package utils

import (
	"errors"
	"reflect"
)

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

func ReverseSlice(data interface{}) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		panic(errors.New("data must be a slice type"))
	}
	valueLen := value.Len()
	for i := 0; i <= int((valueLen-1)/2); i++ {
		reverseIndex := valueLen - 1 - i
		tmp := value.Index(reverseIndex).Interface()
		value.Index(reverseIndex).Set(value.Index(i))
		value.Index(i).Set(reflect.ValueOf(tmp))
	}
}
