package events

import (
	"github.com/sirupsen/logrus"
	"github.com/suutaku/go-vnc/internal/buffer"
	"github.com/suutaku/go-vnc/internal/display"
)

// SetEncodings handles the client set-encodings event.
type SetEncodings struct{}

// Code returns the code.
func (s *SetEncodings) Code() uint8 { return 2 }

// Handle handles the event.
func (s *SetEncodings) Handle(buf *buffer.ReadWriter, d *display.Display) error {
	if err := buf.ReadPadding(1); err != nil {
		return err
	}

	var numEncodings uint16
	if err := buf.Read(&numEncodings); err != nil {
		return err
	}

	encTypes := make([]int32, int(numEncodings))
	for i := 0; i < int(numEncodings); i++ {
		if err := buf.Read(&encTypes[i]); err != nil {
			return err
		}
	}
	encs, pseudo := splitPseudoEncodings(encTypes)
	logrus.Infof("Client encodings: %#v", encs)
	logrus.Infof("Client pseudo-encodings: %#v", pseudo)
	d.SetEncodings(encs, pseudo)

	return nil
}

func splitPseudoEncodings(encs []int32) (encodings, pseudoEncodings []int32) {
	encodings = make([]int32, 0)
	var i int
	for i = 0; i < len(encs); i++ {
		encodings = append(encodings, encs[i])
		if encs[i] == 0 {
			break
		}
	}
	if i == len(encs)-1 {
		return
	}
	pseudoEncodings = encs[i+1:]
	return
}
