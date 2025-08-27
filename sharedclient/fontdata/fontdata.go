package fontdata

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	DEFAULT_FONT_SIZE = 14
	MAX_LINE_WIDTH    = 315
)

var DEFAULT_FONT_FACE = func() *text.GoTextFace {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	fontFace := &text.GoTextFace{
		Source: s,
		Size:   DEFAULT_FONT_SIZE,
	}
	return fontFace
}()

var SMALLER_FONT_FACE = func() *text.GoTextFace {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	fontFace := &text.GoTextFace{
		Source: s,
		Size:   DEFAULT_FONT_SIZE - 4,
	}
	return fontFace
}()

var TITLE_FONT_FACE = func() *text.GoTextFace {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	fontFace := &text.GoTextFace{
		Source: s,
		Size:   DEFAULT_FONT_SIZE + 6,
	}
	return fontFace
}()

var UNLOCK_BODY_FONT_FACE = func() *text.GoTextFace {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	fontFace := &text.GoTextFace{
		Source: s,
		Size:   DEFAULT_FONT_SIZE - 2,
	}
	return fontFace
}()
