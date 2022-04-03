package events

import (
	"github.com/sirupsen/logrus"
	"github.com/suutaku/go-vnc/internal/buffer"
	"github.com/suutaku/go-vnc/internal/display"
	"github.com/suutaku/go-vnc/internal/types"
)

// SetPixelFormat handles the client set-pixel-format event.
type SetPixelFormat struct{}

// Code returns the code.
func (s *SetPixelFormat) Code() uint8 { return 0 }

// Handle handles the event.
func (s *SetPixelFormat) Handle(buf *buffer.ReadWriter, d *display.Display) error {

	if err := buf.ReadPadding(3); err != nil {
		return err
	}

	var pf types.PixelFormat
	if err := buf.ReadInto(&pf); err != nil {
		return err
	}

	logrus.Infof("Client wants pixel format: %#v", pf)
	d.SetPixelFormat(&pf)

	if err := buf.ReadPadding(3); err != nil {
		return err
	}
	return nil
}
