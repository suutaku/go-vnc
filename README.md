# go-vnc
One of the goals of this project was create a pure go implementation of VNC(RFB).
Based on `https://github.com/tinyzimme/gsvnc`, to make project to 'pure go', I removed the `gstream display impl` which using CGO from gsvnc. At now, this project still not a pure go implementaion of RFB due the dependency of `https://github.com/go-vgo/robotgo.git`.

## Install

just:

```bash
go mod tidy
go build ./cmd/vnc
```

## Using package as a library

this package export two modules from `pkg`. One is the `vnc`, a vnc server library and another one `client` for a vnc client library (But not complate yet >.< , I prefer using clients like noVNC).
The server exmaple:

```golang
package main
import (
  "github.com/suutaku/go-vnc/pkg/vnc"
  "github.com/suutaku/go-vnc/pkg/config"
  "context"
)


func main(){
  vncServer := vnc.NewVNC(context.Backgroud(),config.DefaultConfigure)
  vncServer.Start()
}
```

The `config` is a server setting configure. At exmple, we used `DefaultConfigure` and it's value described as:

```golang
var DefaultConfigure = Configure{
	Debug: true,
	Websockify: WebsockifyConf{
		Port: 2225,
		Host: "127.0.0.1",
	},
	Resolution: ResolutionConf{
		Width:  2880,
		Height: 1800,
	},
	DisplayImpl:  display.ProviderScreenCapture,
	AuthType:     []string{"VNCAuth", "TightSecurity"}, // None, VNCAuth, TightSecurity
	EncodingType: []string{"TightPNGEncoding", "RawEncoding", "TightEncoding"},
	EventType:    []string{"KeyEvent", "PointerEvent", "FrameBufferUpdate", "SetPixelFormat", "SetEncodings", "ClientCutText"},
}
```

for more informations just check `pkg/config/config.go`.