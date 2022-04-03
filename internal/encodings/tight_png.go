package encodings

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"log"
	"strconv"

	"github.com/suutaku/go-vnc/internal/types"
	"github.com/suutaku/go-vnc/internal/utils"
)

// TightPNGEncoding implements an Encoding intercace using Tight encoding.
type TightPNGEncoding struct{}

// Code returns the code
func (t *TightPNGEncoding) Code() int32 { return -260 }

// HandleBuffer handles an image sample.
func (t *TightPNGEncoding) HandleBuffer(w io.Writer, f *types.PixelFormat, img *image.RGBA) {
	compressed := new(bytes.Buffer)

	err := png.Encode(compressed, img)
	if err != nil {
		log.Println("[tight-png] Could not encode image frame to png")
		return
	}

	buf := compressed.Bytes()

	i, _ := strconv.ParseInt("01010000", 2, 64) // PNG encoding
	utils.Write(w, uint8(i))

	// Buffer length
	utils.Write(w, computeTightLength(len(buf)))

	// Buffer contents
	utils.Write(w, buf)
}
