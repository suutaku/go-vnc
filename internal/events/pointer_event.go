package events

import (
	"github.com/suutaku/go-vnc/internal/buffer"
	"github.com/suutaku/go-vnc/internal/display"
	"github.com/suutaku/go-vnc/internal/types"
)

// PointerEvent handles pointer events.
type PointerEvent struct{}

// Code returns the code.
func (s *PointerEvent) Code() uint8 { return 5 }

// Handle handles the event.
func (s *PointerEvent) Handle(buf *buffer.ReadWriter, d *display.Display) error {
	var req types.PointerEvent
	if err := buf.ReadInto(&req); err != nil {
		return err
	}
	d.DispatchPointerEvent(&req)
	return nil
}
