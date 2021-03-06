package rfb

import (
	"net"

	"github.com/sirupsen/logrus"
	"github.com/suutaku/go-vnc/internal/buffer"
	"github.com/suutaku/go-vnc/internal/display"
	"github.com/suutaku/go-vnc/internal/events"
)

// Conn represents a client connection.
type Conn struct {
	c       net.Conn
	s       *Server
	buf     *buffer.ReadWriter
	display *display.Display
}

func (s *Server) newConn(c net.Conn) *Conn {
	buf := buffer.NewReadWriteBuffer(c)
	conn := &Conn{
		c:   c,
		s:   s,
		buf: buf,
		display: display.NewDisplay(&display.Opts{
			Width:           s.width,
			Height:          s.height,
			Buffer:          buf,
			DisplayProvider: s.displayProvider,
			GetEncodingFunc: s.GetEncoding,
		}),
	}
	return conn
}

func (c *Conn) serve() {
	defer c.c.Close()

	if err := c.display.Start(); err != nil {
		logrus.Errorf("Error starting display: %s", err)
		return
	}
	defer c.display.Close()

	// Get a map of event handlers for this connection
	eventHandlers := c.s.GetEventHandlerMap()
	defer events.CloseEventHandlers(eventHandlers)

	// handle events
	for {
		cmd, err := c.buf.ReadByte()
		if err != nil {
			logrus.Errorf("Client disconnect: %s", err.Error())
			return
		}
		if hdlr, ok := eventHandlers[cmd]; ok {
			if err := hdlr.Handle(c.buf, c.display); err != nil {
				logrus.Errorf("Error handling cmd %d: %s", cmd, err.Error())
				return
			}
		} else {
			logrus.Warningf("Unsupported command type %d from client\n", int(cmd))
		}
	}
}
