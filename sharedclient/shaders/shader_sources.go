package shaders

import (
	_ "embed" // Don't forget the underscore import
)

//go:embed fog.kage
var fogShaderSrc []byte
