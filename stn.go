package stn

import (
	"fmt"
	"github.com/kbinani/screenshot"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

type Monitor struct {
	Root   int
	imgs   []*image.RGBA
	acc    int
	active bool
}

/**
w: resize weight
h: resize height
v: concat dim
images: list
*/
func Concat(w uint, h uint, v bool, images ...image.Image) *image.RGBA {

	var acc uint = 0
	if w == 0 {
		v = false
	} else if h == 0 {
		v = true
	}

	for i, img := range images {
		rimg := resize.Resize(w, h, img, resize.Bilinear)
		if v { // vertical concat, accumulate height
			acc += uint(rimg.Bounds().Dy())
		} else {
			acc += uint(rimg.Bounds().Dx())
		}
		images[i] = rimg
	}

	if v {
		h = acc
	} else {
		w = acc
	}

	r := image.Rectangle{image.Point{0, 0}, image.Point{int(w), int(h)}}
	rgba := image.NewRGBA(r)

	dx := 0
	dy := 0

	for _, img := range images {

		rec := img.Bounds()
		draw.Draw(rgba, image.Rect(dx, dy, dx+rec.Dx(), dy+rec.Dy()), img, image.Point{0, 0}, draw.Src)
		if v {
			dy += img.Bounds().Dy()
		} else {
			dx += img.Bounds().Dx()
		}
	}

	return rgba
}

func IsImageEqual(a *image.RGBA, b *image.RGBA, thresh int32) bool {
	//ahash := fmt.Sprintf("%x", md5.Sum(a.Pix))
	//bhash := fmt.Sprintf("%x", md5.Sum(b.Pix))
	//println("%s,%s", ahash, bhash)
	cnt, _ := byteDiff(a.Pix, b.Pix)
	fmt.Println("%d %d", cnt, len(a.Pix))
	return cnt <= thresh
}

func byteDiff(bs1, bs2 []byte) (int32, error) {
	// Ensure that we have two non-nil slices with the same length.
	if (bs1 == nil) || (bs2 == nil) {
		return -1, fmt.Errorf("expected a byte slice but got nil")
	}
	if len(bs1) != len(bs2) {
		return -1, fmt.Errorf("mismatched lengths, %d != %d", len(bs1), len(bs2))
	}

	// Populate and return the difference between the two.
	cnt := int32(0)
	for i := range bs1 {
		if bs1[i] != bs2[i] {
			cnt += 1
		}
	}
	return cnt, nil
}

func GetFullScreenShot(w uint, h uint, v bool) *image.RGBA {
	n := screenshot.NumActiveDisplays()
	imgs := []image.Image{}
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		imgs = append(imgs, img)
	}
	cimg := Concat(w, h, v, imgs...)

	return cimg
}

func UpdateScreenshotInfo(file string, info *image.RGBA) (uint, string) {
	f, err := os.OpenFile(file, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		return 1, fmt.Sprint(err)
	}

	err = jpeg.Encode(f, info, nil)
	if err != nil {
		return 2, fmt.Sprint(err)
	}

	return 0, ""
}
