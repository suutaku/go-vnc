package config

// Debug represents if debug logging is enabled. It is mutated at boot.
var Debug = false

type WebsockifyConf struct {
	Host string
	Port int32
}

type TCPConf struct {
	Host string
	Port int32
}

type ResolutionConf struct {
	Width  int32
	Height int32
}

type Configure struct {
	Debug        bool
	TCP          TCPConf
	Resolution   ResolutionConf
	AuthFilePath string
	DisplayImpl  string
	Websockify   WebsockifyConf
	AuthType     []string
	EncodingType []string
	EventType    []string
}
