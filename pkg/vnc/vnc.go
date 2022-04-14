package vnc

import (
	"context"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/suutaku/go-vnc/internal/display"
	"github.com/suutaku/go-vnc/internal/rfb"
	"github.com/suutaku/go-vnc/internal/utils"
	"github.com/suutaku/go-vnc/pkg/config"
)

type VNC struct {
	conf   config.Configure
	server *rfb.Server
}

func NewVNC(ctx context.Context, conf config.Configure) *VNC {
	opts := &rfb.ServerOpts{
		Width:            int(conf.Resolution.Width),
		Height:           int(conf.Resolution.Height),
		DisplayProvider:  display.Provider(conf.DisplayImpl),
		EnabledAuthTypes: configureAuthTypes(conf.AuthType),
		EnabledEncodings: configureEncodings(conf.EncodingType),
		EnabledEvents:    configureEvents(conf.EventType),
		ServerPassword:   conf.Password,
	}

	if authIsEnabled(opts.EnabledAuthTypes, "VNCAuth") {
		if opts.ServerPassword == "" {
			logrus.Info("VNCAuth is enabled and no password provided, generating a server password")
			opts.ServerPassword = utils.RandomString(8)
			logrus.Info("Clients using VNCAuth can connect with the following password: ", opts.ServerPassword)
		}
	}

	return &VNC{
		server: rfb.NewServer(opts),
		conf:   conf,
	}
}

func (vnc *VNC) Start() {
	noTCP := false
	noWebsockify := false
	if vnc.conf.TCP.Host == "" || vnc.conf.TCP.Port == 0 {
		noTCP = true
	}
	if vnc.conf.Websockify.Host == "" || vnc.conf.Websockify.Port == 0 {
		noWebsockify = true
	}
	if noTCP && noWebsockify {
		panic("no listen service! please enable TCP or WS")
	}
	if !noWebsockify {
		go vnc.serveWebsockify()
	}
	if !noTCP {
		// Create a listener
		bindAddr := fmt.Sprintf("%s:%d", vnc.conf.TCP.Host, vnc.conf.TCP.Port)
		l, err := net.Listen("tcp", bindAddr)
		if err != nil {
			panic(err)
		}
		logrus.Info("listening for rfb connections on ", bindAddr)
		vnc.server.Serve(l)
	}
}

func (vnc *VNC) serveWebsockify() {
	wsAddr := fmt.Sprintf("%s:%d", vnc.conf.Websockify.Host, vnc.conf.Websockify.Port)
	l, err := net.Listen("tcp", wsAddr)
	if err != nil {
		panic(err)
	}
	logrus.Info("listening for websockify connections on ", wsAddr)
	vnc.server.ServeWebsockify(l)
}
