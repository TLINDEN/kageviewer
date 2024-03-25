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
	"log"
	"os"
	"runtime/debug"

	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tlinden/yadu"
)

func main() {
	// parse config file and command line parameters, if any
	conf, err := InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	if conf.Showversion {
		fmt.Printf("This is kage-viewer version %s\n", VERSION)
		return
	}

	//  enable  debugging,  if  needed.   We  only  use  log/slog  for
	// debugging, so there's no need to configure it outside debugging
	if conf.Debug {
		logLevel := &slog.LevelVar{}
		// we're using a more verbose logger in debug mode
		buildInfo, _ := debug.ReadBuildInfo()
		opts := &yadu.Options{
			Level:     logLevel,
			AddSource: true,
		}

		logLevel.Set(slog.LevelDebug)

		handler := yadu.NewHandler(os.Stdout, opts)
		debuglogger := slog.New(handler).With(
			slog.Group("program_info",
				slog.Int("pid", os.Getpid()),
				slog.String("go_version", buildInfo.GoVersion),
			),
		)
		slog.SetDefault(debuglogger)
	}

	game := &Game{Conf: conf}
	if err := game.Init(); err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(game.Conf.Width, game.Conf.Height)
	ebiten.SetWindowTitle("Kage shader viewer")
	ebiten.SetWindowResizable(true)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
