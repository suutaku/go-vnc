package encodings

import (
	"bytes"
	"image"
	"io"
	"log"
	"strconv"

	"github.com/pixiv/go-libjpeg/jpeg"

	"github.com/suutaku/go-vnc/internal/types"
	"github.com/suutaku/go-vnc/internal/utils"
)

// TightEncoding implements an Encoding intercace using Tight encoding.
type TightEncoding struct{}

// Code returns the code
func (t *TightEncoding) Code() int32 { return 7 }

// HandleBuffer handles an image sample.
func (t *TightEncoding) HandleBuffer(w io.Writer, f *types.PixelFormat, img *image.RGBA) {

	compressed := new(bytes.Buffer)

	err := jpeg.Encode(compressed, img, &jpeg.EncoderOptions{Quality: 80, OptimizeCoding: true, DCTMethod: jpeg.DCTFloat})
	if err != nil {
		log.Println("[tight-jpeg] Could not encode image frame to jpeg")
		return
	}

	buf := compressed.Bytes()

	i, _ := strconv.ParseInt("10010000", 2, 64) // JPEG encoding
	utils.Write(w, uint8(i))

	// Buffer length
	utils.Write(w, computeTightLength(len(buf)))

	// Buffer contents
	utils.Write(w, buf)
}

func computeTightLength(compressedLen int) (b []byte) {
	out := []byte{byte(compressedLen & 0x7F)}
	if compressedLen > 0x7F {
		out[0] |= 0x80
		out = append(out, byte(compressedLen>>7&0x7F))
		if compressedLen > 0x3FFF {
			out[1] |= 0x80
			out = append(out, byte(compressedLen>>14&0xFF))
		}
	}
	return out
}
