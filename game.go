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

	"github.com/Zyko0/Ebiary/asset"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Conf   *Config
	Images []*asset.LiveAsset[*ebiten.Image]
	Shader *asset.LiveAsset[*ebiten.Shader]
	Cursor []float64
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
		//img, err := LoadImage(image)
		img, err := asset.NewLiveAsset[*ebiten.Image](image)
		if err != nil {
			return fmt.Errorf("failed to load image %s: %s", image, err)
		}

		game.Images = append(game.Images, img)
	}

	shader, err := asset.NewLiveAsset[*ebiten.Shader](game.Conf.Shader)
	if err != nil {
		return fmt.Errorf("failed to load shader %s: %s", game.Conf.Shader, err)
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
	for _, image := range game.Images {
		if image.Error() != nil {
			fmt.Println("warn: image reloading error:", image.Error())
		}
	}

	if game.Shader.Error() != nil {
		fmt.Println("warn: shader reloading error:", game.Shader.Error())
	}

	if game.CheckInput() {
		slog.Debug("Key pressed",
			game.Conf.Flag, game.Flag,
			game.Conf.Slider, game.Slider,
			game.Conf.Ticks, fmt.Sprintf("%.02f", float64(game.Ticks)/60),
			game.Conf.Mouse, fmt.Sprintf("%.02f, %.02f", game.Cursor[0], game.Cursor[1]),
		)
	}

	mousex, mousey := ebiten.CursorPosition()
	game.Cursor = []float64{float64(mousex), float64(mousey)}

	game.Ticks++

	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawRectShaderOptions{}

	op.Uniforms = map[string]any{
		game.Conf.Flag:   game.Flag,
		game.Conf.Slider: game.Slider,
		game.Conf.Ticks:  float64(game.Ticks) / 60,
		game.Conf.Mouse:  game.Cursor,
	}

	for idx, image := range game.Images {
		op.Images[idx] = image.Value()
	}

	op.GeoM.Translate(float64(game.Conf.X), float64(game.Conf.Y))

	screen.DrawRectShader(game.Conf.Width, game.Conf.Height, game.Shader.Value(), op)
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return game.Conf.Width, game.Conf.Height
}
