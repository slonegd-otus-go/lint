package main

type Example1 struct {
	Foo1 string `json:"foo"`
	Bar1 int    `bson:"bar"`
}

type Example2 struct {
	Foo2 string `json:"foo"`
	Bar2 int
}
