package main

import "fmt"

type Example1 struct {
	Foo1 string `json:"foo" bson:"foo"`
	Bar1 int    `bson:"bar"`
}

type Example2 struct {
	Foo2 string `json:"foo"`
	Bar2 int
}

func ExampleFunc(v Example2) string {
	return fmt.Sprintf("%v", v)
}
