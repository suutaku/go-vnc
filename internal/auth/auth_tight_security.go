package auth

import (
	"bytes"
	"fmt"
	"io"

	"github.com/suutaku/go-vnc/internal/buffer"
	"github.com/suutaku/go-vnc/internal/types"
	"github.com/suutaku/go-vnc/internal/utils"
)

// TightTunnelCapabilities represents TightSecurity tunnel capabilities.
var TightTunnelCapabilities = []types.TightCapability{
	{Code: 0, Vendor: "TGHT", Signature: "NOTUNNEL"},
}

// TightAuthCapabilities represents TightSecurity auth capabilities.
var TightAuthCapabilities = []types.TightCapability{
	{Code: 1, Vendor: "STDV", Signature: "NOAUTH__"},
	{Code: 2, Vendor: "STDV", Signature: "VNCAUTH_"},
}

// TightServerMessages represents supported tight server messages.
var TightServerMessages = []types.TightCapability{}

// TightClientMessages represents supported tight client messages.
var TightClientMessages = []types.TightCapability{}

// TightEncodingCapabilities represents TightSecurity encoding capabilities.
var TightEncodingCapabilities = []types.TightCapability{
	{Code: 0, Vendor: "STDV", Signature: "RAW_____"},
	{Code: 1, Vendor: "STDV", Signature: "COPYRECT"},
	{Code: 7, Vendor: "TGHT", Signature: "TIGHT___"},
}

// TightSecurity implements Tight security.
// https://github.com/rfbproto/rfbproto/blob/master/rfbproto.rst#tight-security-type
type TightSecurity struct {
	AuthGetter func(code uint8) Type
}

// Code returns the code.
func (t *TightSecurity) Code() uint8 { return 16 }

// Negotiate will negotiate tight security.
func (t *TightSecurity) Negotiate(rw *buffer.ReadWriter) error {
	if err := t.negotiateTightTunnel(rw); err != nil {
		return err
	}
	return t.negotiateTightAuth(rw)
}

// Response negotiate from client
func (a *TightSecurity) Response(rw *buffer.ReadWriter) error {
	return nil
}

// ExtendServerInit signals to the rfb server that we extend the ServerInit message.
func (t *TightSecurity) ExtendServerInit(buf io.Writer) {
	utils.Write(buf, uint16(len(TightServerMessages)))
	utils.Write(buf, uint16(len(TightClientMessages)))
	utils.Write(buf, uint16(len(TightEncodingCapabilities)))
	utils.Write(buf, uint8(0)) // padding
	utils.Write(buf, uint8(0)) // padding
	for _, cap := range TightServerMessages {
		utils.PackStruct(buf, &cap)
	}
	for _, cap := range TightClientMessages {
		utils.PackStruct(buf, &cap)
	}
	for _, cap := range TightEncodingCapabilities {
		utils.PackStruct(buf, &cap)
	}
}

func (t *TightSecurity) negotiateTightTunnel(rw *buffer.ReadWriter) error {
	// Write the supported tunnel capabilities to the client
	buf := new(bytes.Buffer)
	utils.Write(buf, uint32(len(TightTunnelCapabilities)))
	for _, cap := range TightTunnelCapabilities {
		utils.PackStruct(buf, &cap)
	}
	rw.Dispatch(buf.Bytes())

	// get the desired tunnel type from the client
	var tun int32
	rw.Read(&tun)

	// We only support no tunneling for now, client should know
	// better
	if tun != 0 {
		return fmt.Errorf("client requested unsupported tunnel type: %d", tun)
	}

	return nil
}

func (t *TightSecurity) negotiateTightAuth(rw *buffer.ReadWriter) error {
	buf := new(bytes.Buffer)
	caps := t.getEnabledAuthCaps()
	utils.Write(buf, uint32(len(caps)))
	for _, cap := range caps {
		utils.PackStruct(buf, &cap)
	}
	rw.Dispatch(buf.Bytes())

	// get the desired auth type, should be none
	var auth int32
	rw.Read(&auth)

	authType := t.AuthGetter(uint8(auth))
	if authType == nil {
		return fmt.Errorf("client requested unsupported tight auth type: %d", auth)
	}
	return authType.Negotiate(rw)
}

func (t *TightSecurity) getEnabledAuthCaps() []types.TightCapability {
	enabledCaps := make([]types.TightCapability, 0)
	for _, cap := range TightAuthCapabilities {
		if t := t.AuthGetter(uint8(cap.Code)); t != nil {
			enabledCaps = append(enabledCaps, cap)
		}
	}
	return enabledCaps
}
