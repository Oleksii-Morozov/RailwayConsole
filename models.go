package main

type Station struct {
	Id   int64
	City string
	Name string
}

type StationCreate struct {
	City string
	Name string
}

type Train struct {
	Id       int64
	Code     string
	Capacity int
}

type TrainCreate struct {
	Code     string
	Capacity int
}
