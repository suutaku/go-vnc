package client

import (
	"fmt"

	"github.com/suutaku/go-vnc/internal/buffer"
)

func (cl *Client) responseAuthNegotiate(ver string, buf *buffer.ReadWriter) error {
	var numSecuType uint8
	err := buf.Read(&numSecuType)
	if err != nil {
		return fmt.Errorf("unknow server auth type")
	}
	if numSecuType == 0 {
		return fmt.Errorf("no security types")
	}
	securityTypes := make([]uint8, numSecuType)
	if err := buf.Read(&securityTypes); err != nil {
		return err
	}
	for _, v := range securityTypes {
		if cl.authType.Code() == v {
			buf.Dispatch([]byte{cl.authType.Code()})

			return nil
		}
	}
	return fmt.Errorf("unsuported server auth type")
}
