package main

func Foo1() ([]int, error) {
	return nil, nil
}

func Foo2() []int {
	return nil
}

func Foo3() {}

func Foo4() error {
	return nil
}

func Foo5() ([]int, error) {
	return Foo1()
}

func Foo6() ([]int, error) {
	var res []int
	return res, nil
}

func Foo7() ([]int, error) {
	err := Foo4()
	if err != nil {
		return nil, err
	} else {
		return nil, err
	}
	tmp := []int{}
	for _, _ = range tmp {
		return nil, err
	}
	return nil, nil
}

// func Foo2() ([]int, error) {
// 	return make([]int, 0), nil
// }

// func Foo3() ([]int, error) {
// 	res, err := Foo2()
// 	return res, err
// }

// func Foo4() ([]int, error) {
// 	return Foo1()
// }
