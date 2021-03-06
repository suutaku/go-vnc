package display

import (
	"image"

	"github.com/sirupsen/logrus"
)

// A Display is an interface that can be implemented by different types of frame sources.
//
// It is the responsibility of these providers to tap into some video/image source and
// make frames available for procesesing.
type IDisplay interface {
	// Start should take care of any requirements for starting a feed to the frame buffer.
	Start(width, height int) error
	// PullFrame should return a queued frame for processing.
	PullFrame() *image.RGBA
	// Close should stop any background processes from running.
	Close() error
}

// Provider is an enum used for selecting a display provider.
type Provider string

// Provider options.
const (
	ProviderGstreamer  = "gstreamer"
	ProviderScreenShot = "screenshot"
)

// GetDisplayProvider returns the provider to use for the given RFB connection.
func GetDisplayProvider(p Provider) IDisplay {
	switch p {
	case ProviderGstreamer:
		logrus.Warn("gstreamer impl was removed")
		// return &Gstreamer{}
		return nil
	case ProviderScreenShot:
		logrus.Info("screen shot provider")
		return &ScreenShot{}
	}
	return nil
}
