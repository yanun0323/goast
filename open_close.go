package goast

type OpenClose int8

const (
	OpenCloseNone OpenClose = iota
	OpenCloseOpen
	OpenCloseClose
)
