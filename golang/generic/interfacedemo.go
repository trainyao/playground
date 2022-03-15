package main

import (
	"log"
)

func caller() {
	var intParams = 1
	testSum(intParams)
	testSum("string")
	testSum(3.33)
	testSum([]byte{1, 2, 3})
}

func testSum(param ...interface{}) interface{} {
	switch param[0].(type) {
	case int:
		return sumInt(param)
	case string:
		return concatString(param)
	case float64:
		// TODO
	default:
		log.Fatalln("黑人问号，不知道怎么进行累加")
	}
	return nil
}

func sumInt(ints []interface{}) (res int) {
	for _, i := range ints {
		res += i.(int)
	}
	return
}

func concatString(strings []interface{}) (res string) {
	// TODO
	return ""
}
