package display

import (
	"image"

	"github.com/sirupsen/logrus"
	"github.com/suutaku/screenshot/pkg/screenshot"
	"golang.org/x/image/draw"
)

type ScreenShot struct {
	frameQueue   chan *image.RGBA
	stopCh       chan struct{}
	screenshoter *screenshot.Screenshot
}

func NewScreenShot() *ScreenShot {
	return &ScreenShot{}
}

// Start should take care of any requirements for starting a feed to the frame buffer.
func (ss *ScreenShot) Start(width, height int) error {

	ss.screenshoter = screenshot.NewScreenshot(0, 0, 0, 0)
	ss.frameQueue = make(chan *image.RGBA, 2)
	ss.stopCh = make(chan struct{})
	dr := image.Rect(0, 0, width, height)
	newImg := image.NewRGBA(dr)
	go func() {
		logrus.Info("display [ScreenShot] start")
		for {
			img, err := ss.screenshoter.Capture()
			if err != nil {
				logrus.Error(err)
				ss.Close()
				return
			}
			if img.Bounds().Max.X > width || img.Bounds().Max.Y > height {
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
	if ss.screenshoter != nil {
		ss.screenshoter.Close()
	}
	return nil
}
