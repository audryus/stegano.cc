package control

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io"
)

func Embed(message, filename string, file io.Reader) (*bytes.Buffer, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	// Prepend 4-byte length (big-endian) so decoder knows how many bytes to read
	payloadLen := uint32(len(message))
	payload := make([]byte, 4+len(message))
	binary.BigEndian.PutUint32(payload[:4], payloadLen)
	copy(payload[4:], message)

	// Embed payload bits across the whole image and get output image
	outImg := embedPayloadIntoImage(img, payload)

	var buf bytes.Buffer
	if err := png.Encode(&buf, outImg); err != nil {
		return nil, fmt.Errorf("failed to encode output png: %w", err)
	}
	return &buf, nil
}

// embedPayloadIntoImage embeds payload bits LSB-first across R,G,B channels and
// returns a new RGBA image with modified LSBs. Non-embedded pixels are copied.
func embedPayloadIntoImage(src image.Image, payload []byte) *image.RGBA {
	bounds := src.Bounds()
	out := image.NewRGBA(bounds)

	bitIdx := 0
	totalBits := len(payload) * 8

	nextBit := func() (uint8, bool) {
		if bitIdx >= totalBits {
			return 0, false
		}
		byteIdx := bitIdx / 8
		bitInByte := 7 - (bitIdx % 8)
		b := (payload[byteIdx] >> bitInByte) & 1
		bitIdx++
		return b, true
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)

			if bit, ok := nextBit(); ok {
				r8 = (r8 & 0xFE) | bit
			}
			if bit, ok := nextBit(); ok {
				g8 = (g8 & 0xFE) | bit
			}
			if bit, ok := nextBit(); ok {
				b8 = (b8 & 0xFE) | bit
			}

			out.Set(x, y, color.RGBA{r8, g8, b8, a8})
		}
	}

	return out
}
