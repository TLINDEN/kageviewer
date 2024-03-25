/*
Copyright © 2024 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Conf   *Config
	Images []*ebiten.Image
	Shader *ebiten.Shader
	Ticks  int
	Slider float64
	Flag   int
}

func LoadImage(name string) (*ebiten.Image, error) {
	fd, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	img, _, err := image.Decode(fd)
	if err != nil {
		return nil, fmt.Errorf("failed to load image %s: %s", name, err)
	}

	return ebiten.NewImageFromImage(img), nil
}

func (game *Game) Init() error {
	for _, image := range game.Conf.Image {
		slog.Debug("Loading images", "image", image)
		img, err := LoadImage(image)
		if err != nil {
			return err
		}

		game.Images = append(game.Images, img)
	}

	data, err := os.ReadFile(game.Conf.Shader)
	if err != nil {
		return fmt.Errorf("failed to load shader %s: %s", game.Conf.Shader, err)
	}

	shader, err := ebiten.NewShader(data)
	if err != nil {
		return fmt.Errorf("failed to create new shader %s: %s", game.Conf.Shader, err)
	}

	game.Shader = shader

	return nil
}

func (g *Game) CheckInput() bool {
	pressed := false

	switch {
	case inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0):
		fallthrough
	case inpututil.IsKeyJustPressed(ebiten.KeySpace):
		g.Toggle()
		pressed = true
	case inpututil.IsKeyJustPressed(ebiten.KeyUp):
		g.Up()
		pressed = true
	case inpututil.IsKeyJustPressed(ebiten.KeyDown):
		g.Down()
		pressed = true
	case inpututil.IsKeyJustPressed(ebiten.KeyQ):
		os.Exit(0)
	}

	return pressed
}

func (g *Game) Toggle() {
	g.Flag = 1 - g.Flag
}

func (g *Game) Up() {
	g.Slider += 0.1
	if g.Slider > 1.0 {
		g.Slider = 1.0
	}
}

func (g *Game) Down() {
	g.Slider -= 0.1

	if g.Slider < 0 {
		g.Slider = 0.0
	}
}

func (game *Game) Update() error {
	if game.CheckInput() {
		slog.Debug("Key pressed", "Slider", game.Slider, "Flag", game.Flag)
	}

	game.Ticks++

	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawRectShaderOptions{}

	mousex, mousey := ebiten.CursorPosition()

	op.Uniforms = map[string]any{
		"Flag":   game.Flag,
		"Slider": game.Slider,
		"Ticks":  game.Ticks,
		"Mouse":  []float64{float64(mousex), float64(mousey)},
	}

	copy(op.Images[:3], game.Images)

	op.GeoM.Translate(float64(game.Conf.X), float64(game.Conf.Y))

	screen.DrawRectShader(game.Conf.Width, game.Conf.Height, game.Shader, op)
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return game.Conf.Width, game.Conf.Height
}
