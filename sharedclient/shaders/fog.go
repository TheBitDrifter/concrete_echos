package shaders

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var FogShader *ebiten.Shader = func() *ebiten.Shader {
	fogShader, err := ebiten.NewShader(fogShaderSrc)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(fogShader, "fs")
	return fogShader
}()
