package client

import (
	"net"

	"github.com/suutaku/go-vnc/internal/auth"
	"github.com/suutaku/go-vnc/internal/buffer"
	"github.com/suutaku/go-vnc/internal/types"
	"github.com/suutaku/go-vnc/internal/version"
)

type Client struct {
	conn           net.Conn
	buf            *buffer.ReadWriter
	authType       auth.Type
	version        string
	frameBufWidth  uint16
	frameBufHeight uint16
	pixelFormat    *types.PixelFormat
	serviceName    string
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:    conn,
		buf:     buffer.NewReadWriteBuffer(conn),
		version: version.V8,
	}
}

func (cl *Client) Handshake() error {
	ver, err := version.ResponseProtocolVersion(cl.buf)
	if err != nil {
		return err
	}
	err = cl.responseAuthNegotiate(ver, cl.buf)
	if err != nil {
		return err
	}
	err = cl.authType.Response(cl.buf)
	if err != nil {
		return err
	}

	//client init
	var sharedFlag uint8 = 1
	cl.buf.Dispatch([]byte{sharedFlag})

	// server init
	if err := cl.buf.Read(&cl.frameBufWidth); err != nil {
		return err
	}
	if err := cl.buf.Read(&cl.frameBufHeight); err != nil {
		return err
	}
	if err := cl.buf.Read(cl.pixelFormat); err != nil {
		return err
	}
	var pdValue uint8
	if err := cl.buf.Read(&pdValue); err != nil {
		return err
	}
	if err := cl.buf.Read(&pdValue); err != nil {
		return err
	}
	if err := cl.buf.Read(&pdValue); err != nil {
		return err
	}
	var nameLength uint32
	if err := cl.buf.Read(&nameLength); err != nil {
		return err
	}
	if err := cl.buf.Read(&cl.serviceName); err != nil {
		return err
	}

	return nil
}
