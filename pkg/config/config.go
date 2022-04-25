package config

import (
	"github.com/suutaku/go-vnc/internal/display"
	"github.com/suutaku/go-vnc/internal/utils"
)

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
	Password     string
}

var DefaultConfigure = Configure{
	Debug: true,
	Websockify: WebsockifyConf{
		Port: 2225,
		Host: "127.0.0.1",
	},
	Password:     utils.RandomString(8),
	DisplayImpl:  display.ProviderScreenShot,
	AuthType:     []string{"VNCAuth", "TightSecurity"}, // None, VNCAuth, TightSecurity
	EncodingType: []string{"TightPNGEncoding", "RawEncoding", "TightEncoding"},
	EventType:    []string{"KeyEvent", "PointerEvent", "FrameBufferUpdate", "SetPixelFormat", "SetEncodings", "ClientCutText"},
}
