package grayt

import (
	"fmt"
	"image"
)

func traceImage(fl flags, a *API, completed *int64) image.Image {
	pxHigh := fl.pxWide * a.aspectRatio[1] / a.aspectRatio[0]
	fmt.Println(fl.pxWide, pxHigh)
	return image.NewNRGBA(image.Rect(0, 0, fl.pxWide, pxHigh))
}
