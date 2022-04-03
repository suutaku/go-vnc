package events

import (
	"github.com/suutaku/go-vnc/internal/buffer"
	"github.com/suutaku/go-vnc/internal/display"
	"github.com/suutaku/go-vnc/internal/types"
)

// ClientCutText handles new text in the client's cut buffer.
type ClientCutText struct{}

// Code returns the code.
func (c *ClientCutText) Code() uint8 { return 6 }

// Handle handles the event.
func (c *ClientCutText) Handle(buf *buffer.ReadWriter, d *display.Display) error {
	var req types.ClientCutText

	buf.ReadPadding(3)

	if err := buf.Read(&req.Length); err != nil {
		return err
	}

	req.Text = make([]byte, req.Length)

	if err := buf.Read(&req.Text); err != nil {
		return err
	}

	d.DispatchClientCutText(&req)
	return nil
}
