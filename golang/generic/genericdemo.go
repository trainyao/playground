package main

import "fmt"

type underlayInt int

func genericSum[T string | ~int](params ...T) (res T) {
	for _, i := range params {
		res += i
	}
	return
}

type myInterface interface {
	comparable
	string | ~int
	String() string
}

type myInterfaceImpl int

func (u myInterfaceImpl) String() string {
	return "underlayInt"
}

func genericSumWithInterface[T myInterface](params ...T) (res T) {
	for _, i := range params {
		res += i
	}
	return
}

func genericCaller() {
	fmt.Println("--in genericCaller--")
	fmt.Println(genericSum(1, 2, 3))
	fmt.Println(genericSum("string", "string2"))
	fmt.Println(genericSum(underlayInt(1), underlayInt(2)))
}

func genericWithInterfaceCaller() {
	fmt.Println("--in genericWithInterfaceCaller--")
	//fmt.Println(genericSumWithInterface(1, 2, 3))             // invalid, int does not implement myInterface
	//fmt.Println(genericSumWithInterface("string", "string2")) // invalid, string does not implement myInterface
	fmt.Println(genericSumWithInterface(myInterfaceImpl(1), myInterfaceImpl(2)))
}

func main() {
	genericCaller()
	genericWithInterfaceCaller()
}
