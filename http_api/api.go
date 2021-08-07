package http_api

var ModuleHandleContainer = make(map[string]ModuleHandles)

type ModuleHandles struct {
	Init    InitFunc
	Handles map[string]HandleFunc
}
