package display

import (
	"time"

	"github.com/sirupsen/logrus"
)

func (d *Display) handleKeyEvents() {
	for {
		select {
		case ev, ok := <-d.keyEvQueue:
			if !ok {
				// Client disconnected.
				return
			}
			logrus.Debug("Got key event: ", ev)
			if ev.IsDown() {
				d.appendDownKeyIfMissing(ev.Key)
				d.dispatchDownKeys()
			} else {
				d.removeDownKey(ev.Key)
			}
		}
	}
}

func (d *Display) handlePointerEvents() {
	for {
		select {
		case ev, ok := <-d.ptrEvQueue:
			if !ok {
				// Client disconnected.
				return
			}
			logrus.Debug("Got pointer event: ", ev)
			d.servePointerEvent(ev)
		}
	}
}

func (d *Display) handleFrameBufferEvents() {
	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		// Framebuffer update requests
		case ur, ok := <-d.fbReqQueue:
			if !ok {
				// Client disconnected.
				return
			}
			logrus.Debug("Handling framebuffer update request")
			d.pushFrame(ur)

		// Send a frame update anyway if there are no updates on the queue
		case <-ticker.C:
			logrus.Debug("Pushing latest frame to client")
			last := d.GetLastImage()
			d.pushImage(last)
		}
	}
}

func (d *Display) handleCutTextEvents() {
	for {
		select {
		case ev, ok := <-d.cutTxtEvsQ:
			if !ok {
				// Client disconnected.
				return
			}
			logrus.Debug("Got cut-text event: ", ev)
			d.syncToClipboard(ev)
		}
	}
}

func (d *Display) watchChannels() {
	go d.handleKeyEvents()
	go d.handlePointerEvents()
	go d.handleFrameBufferEvents()
	go d.handleCutTextEvents()
}
