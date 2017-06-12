package utils

import (
	"errors"
	"math"
)

var BLACKCOLOR uint32 = 0xff000000
var WHITECOLOR uint32 = 0xffffffff
var MIN_ALPHA_SEARCH_MAX_ITERATIONS uint32 = 10
var MIN_ALPHA_SEARCH_PRECISION uint32 = 10
var MIN_CONTRAST_TITLE_TEXT float64 = 3.0
var MIN_CONTRAST_BODY_TEXT float64 = 4.5

func RedColor(c uint32) uint32 {
	r := (uint32(c) >> 16) & 0xFF
	return r
}

func BlueColor(c uint32) uint32 {
	b := uint32(c) & 0xFF
	return b
}

func GreenColor(c uint32) uint32 {
	g := (uint32(c) >> 8) & 0xFF
	return g
}

func AlphaColor(c uint32) uint32 {
	a := (uint32(c) >> 24) & 0xFF
	return a
}

func compositeColors(fg uint32, bg uint32) uint32 {
	var afg, abg, rfg, rbg, gfg, gbg, bfg, bbg float64
	afg = float64(AlphaColor(fg))
	abg = float64(AlphaColor(bg))
	rfg = float64(RedColor(fg))
	rbg = float64(RedColor(bg))
	bfg = float64(BlueColor(fg))
	bbg = float64(BlueColor(bg))
	gfg = float64(GreenColor(fg))
	gbg = float64(GreenColor(bg))
	var alpha1 float64 = afg / 255
	var alpha2 float64 = abg / 255
	var a float64 = (alpha1 + alpha2) * (1 - alpha1)
	var r float64 = (rfg * alpha1) + (rbg * alpha2 * (1 - alpha1))
	var g float64 = (gfg * alpha1) + (gbg * alpha2 * (1 - alpha1))
	var b float64 = (bfg * alpha1) + (bbg * alpha2 * (1 - alpha1))
	return (uint32(a)&0xFF)<<24 | (uint32(r)&0xFF)<<16 | (uint32(g)&0xFF)<<8 | (uint32(b) & 0xFF)
}

func modifyAlpha(color uint32, alpha uint32) uint32 {
	return (color & 0x00ffffff) | (alpha << 24)
}

func calculateContrast(foreground uint32, background uint32) (error, *float64) {
	if AlphaColor(background) != 255 {
		return errors.New("background can not be translucent"), nil
	}
	if AlphaColor(foreground) < 255 {
		foreground = compositeColors(foreground, background)
	}
	var luminance1 float64 = calculateLuminance(foreground) + 0.05
	var luminance2 float64 = calculateLuminance(background) + 0.05
	var contrast *float64 = new(float64)
	*contrast = float64(math.Max(luminance1, luminance2) / math.Min(luminance1, luminance2))
	return nil, contrast
}

func calculateLuminance(color uint32) float64 {
	var red float64 = float64(RedColor(color)) / 255.0
	if red < 0.03928 {
		red = red / 12.92
	} else {
		red = math.Pow((red+0.055)/1.055, 2.4)
	}
	var green float64 = float64(GreenColor(color)) / 255.0
	if green < 0.03928 {
		green = green / 12.92
	} else {
		green = math.Pow((green+0.055)/1.055, 2.4)
	}
	var blue float64 = float64(BlueColor(color)) / 255.0

	if blue < 0.03928 {
		blue = blue / 12.92

	} else {
		blue = math.Pow((blue+0.055)/1.055, 2.4)
	}
	return (0.2126 * red) + (0.7152 * green) + (0.0722 * blue)
}

func findMinimumAlpha(foreground uint32, background uint32, minContrastRatio float64) (error, *uint32) {
	if AlphaColor(background) != 255 {
		return errors.New("background can not be translucent"), nil
	}
	var testForeground uint32 = modifyAlpha(foreground, 255)
	err, testRatio := calculateContrast(testForeground, background)
	if err != nil {
		return err, nil
	}
	if *testRatio < minContrastRatio {
		return errors.New("Fully opaque foreground does not have sufficient contrast"), nil
	}
	var numIterations uint32 = 0
	var minAlpha uint32 = 0
	var maxAlpha uint32 = 255
	for numIterations <= MIN_ALPHA_SEARCH_MAX_ITERATIONS && (maxAlpha-minAlpha) > MIN_ALPHA_SEARCH_PRECISION {
		var testAlpha uint32 = (minAlpha + maxAlpha) / 2
		testForeground = modifyAlpha(foreground, testAlpha)
		err, testRatio = calculateContrast(testForeground, background)
		if err != nil {
			return err, nil
		}
		if *testRatio < minContrastRatio {
			minAlpha = testAlpha
		} else {
			maxAlpha = testAlpha
		}
		numIterations++
	}
	return nil, &maxAlpha
}

func getTextColorForBackground(backgroundColor uint32, minContrastRatio float64) (error, uint32) {
	err, whiteMinAlpha := findMinimumAlpha(WHITECOLOR, backgroundColor, minContrastRatio)
	if err != nil {
		return err, WHITECOLOR
	}
	if *whiteMinAlpha >= 0 {
		return nil, modifyAlpha(WHITECOLOR, *whiteMinAlpha)
	}
	err, blackMinAlpha := findMinimumAlpha(BLACKCOLOR, backgroundColor, minContrastRatio)
	if err != nil {
		return err, BLACKCOLOR
	}
	if *blackMinAlpha >= 0 {
		return nil, modifyAlpha(BLACKCOLOR, *blackMinAlpha)
	}
	// This should not happen!
	return nil, BLACKCOLOR
}
