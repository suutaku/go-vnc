package display

import (
	"fmt"
	"image"

	"github.com/kbinani/screenshot"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/draw"
)

type ScreenShot struct {
	frameQueue chan *image.RGBA
	stopCh     chan struct{}
}

func NewScreenShot() *ScreenShot {
	return &ScreenShot{}
}

// Start should take care of any requirements for starting a feed to the frame buffer.
func (ss *ScreenShot) Start(width, height int) error {
	n := screenshot.NumActiveDisplays()
	if n < 1 {
		return fmt.Errorf("no display found")
	}
	ss.frameQueue = make(chan *image.RGBA, 2)
	ss.stopCh = make(chan struct{})
	go func() {
		logrus.Info("display [ScreenShot] start")
		bounds := screenshot.GetDisplayBounds(0)
		dr := image.Rect(0, 0, width, height)
		newImg := image.NewRGBA(dr)
		for {
			img, err := screenshot.CaptureRect(bounds)
			if err != nil {
				logrus.Error(err)
				ss.Close()
				return
			}
			if bounds.Max.X > width || bounds.Max.Y > height {
				draw.BiLinear.Scale(newImg, dr, img, img.Bounds(), draw.Over, nil)
				img = newImg
			}
			select {
			case <-ss.stopCh:
				logrus.Debug("Received event on stop channel, stopping screen capture")
				return
			case ss.frameQueue <- img:
			default:
				// pop the oldest item off the queue
				// and let the next sample try to get in
				logrus.Debug("Client is behind on frames, forcing oldest one off the queue")
				// select {
				// case <-s.frameQueue:
				// }
				<-ss.frameQueue
			}
		}
	}()
	return nil
}

// PullFrame should return a queued frame for processing.
func (ss *ScreenShot) PullFrame() *image.RGBA {
	return <-ss.frameQueue
}

// Close should stop any background processes from running.
func (ss *ScreenShot) Close() error {
	ss.stopCh <- struct{}{}
	return nil
}
