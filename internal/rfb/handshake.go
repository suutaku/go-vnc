package rfb

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/suutaku/go-vnc/internal/auth"
	"github.com/suutaku/go-vnc/internal/buffer"
	"github.com/suutaku/go-vnc/internal/utils"
	"github.com/suutaku/go-vnc/internal/version"
)

func (c *Conn) doHandshake() error {
	ver, err := version.NegotiateProtocolVersion(c.buf)
	if err != nil {
		return err
	}

	var authType auth.Type
	if authType, err = c.negotiateAuth(ver, c.buf); err != nil {
		return err
	}

	logrus.Info("Reading client init")

	// ClientInit
	if _, err := c.buf.ReadByte(); err != nil {
		return err
	}

	logrus.Info("Sending server init")
	format := c.display.GetPixelFormat()

	// 6.3.2. ServerInit
	width, height := c.display.GetDimensions()
	logrus.Debug("W/H", width, height)
	buf := new(bytes.Buffer)
	utils.Write(buf, uint16(width))
	utils.Write(buf, uint16(height))
	utils.PackStruct(buf, format)
	utils.Write(buf, uint8(0)) // pad1
	utils.Write(buf, uint8(0)) // pad2
	utils.Write(buf, uint8(0)) // pad3
	serverName := "go-vnc"
	utils.Write(buf, int32(len(serverName)))
	utils.Write(buf, []byte(serverName))

	// Chcek if we are extending server init. This is only applicable to TightSecurity.
	if extender, ok := authType.(interface{ ExtendServerInit(io.Writer) }); ok {
		extender.ExtendServerInit(buf)
	}
	c.buf.Dispatch(buf.Bytes())
	return nil
}

const (
	statusOK     = 0
	statusFailed = 1
)

// NegotiateAuth wil negotiate authentication on the given connection, for the
// given version.
func (c *Conn) negotiateAuth(ver string, rw *buffer.ReadWriter) (auth.Type, error) {
	buf := new(bytes.Buffer)

	logrus.Info("Negotiating security")

	utils.Write(buf, uint8(len(c.s.enabledAuthTypes)))
	for _, t := range c.s.enabledAuthTypes {
		utils.Write(buf, t.Code())
	}
	rw.Dispatch(buf.Bytes())
	wanted, err := rw.ReadByte()
	if err != nil {
		return nil, err
	}
	if !c.s.AuthIsSupported(wanted) {
		return nil, fmt.Errorf("client wanted unsupported auth type %d", int(wanted))
	}

	authType := c.s.GetAuth(wanted)
	logrus.Info("Using security: ", reflect.TypeOf(authType).Elem().Name())

	if err := authType.Negotiate(rw); err != nil {
		logrus.Error("Authentication failed")
		buf = new(bytes.Buffer)
		utils.Write(buf, uint32(statusFailed))
		rw.Dispatch(buf.Bytes())
		return nil, err
	}

	if ver >= version.V8 {
		// 6.1.3. SecurityResult
		buf = new(bytes.Buffer)
		utils.Write(buf, uint32(statusOK))
		rw.Dispatch(buf.Bytes())
	}

	return authType, nil
}
