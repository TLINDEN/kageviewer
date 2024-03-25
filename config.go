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
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

const (
	VERSION string = "0.0.1"
	Usage   string = `This is kage-viewer, a shader viewer.

Usage: kage-viewer [-vd] [-c <config file>] [-g geom] [-p geom] \
       -i <image0.png> -i <image1.png> -s <shader.kage>

Options:
-c --config   <toml file>     Config file to use (optional)
-i --image    <png file>      Image to load (multiple times allowed, up to 4)
-s --shader   <kage file>     Shader to run
-g --geometry <WIDTHxHEIGHT>  Window size
-p --position <XxY>           Position of image0
-d --debug    Show debugging output
-v --version  Show program version
`
)

type Config struct {
	Showversion bool     `koanf:"version"`  // -v
	Debug       bool     `koanf:"debug"`    // -d
	Config      string   `koanf:"config"`   // -c
	Image       []string `koanf:"image"`    // -i
	Shader      string   `koanf:"shader"`   // -s
	Geo         string   `koanf:"geometry"` // -g
	Posision    string   `koanf:"position"` // -p

	X, Y, Width, Height int // feed from -g + -p
}

func InitConfig() (*Config, error) {
	var kloader = koanf.New(".")

	// Load default values using the confmap provider.
	/* not needed yet
	   if err := kloader.Load(confmap.Provider(map[string]interface{}{
	   }, "."), nil); err != nil {
	           return nil, fmt.Errorf("failed to load default values into koanf: %w", err)
	   }
	*/

	// setup custom usage
	flagset := flag.NewFlagSet("config", flag.ContinueOnError)
	flagset.Usage = func() {
		fmt.Println(Usage)
		os.Exit(0)
	}

	// parse commandline flags
	flagset.BoolP("version", "v", false, "show program version")
	flagset.BoolP("debug", "d", false, "enable debug output")
	flagset.StringP("config", "c", "", "config file")
	flagset.StringP("geometry", "g", "256x256", "window geometry")
	flagset.StringP("position", "p", "0x0", "position of shader")
	flagset.StringArrayP("image", "i", nil, "image file")
	flagset.StringP("shader", "s", "", "shader file")

	if err := flagset.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse program arguments: %w", err)
	}

	// generate a  list of config files to try  to load, including the
	// one provided via -c, if any
	var configfiles []string

	configfile, _ := flagset.GetString("config")
	home, _ := os.UserHomeDir()

	if configfile != "" {
		configfiles = []string{configfile}
	} else {
		configfiles = []string{
			"/etc/kage-viewer.conf", "/usr/local/etc/kage-viewer.conf", // unix variants
			filepath.Join(home, ".config", "kage-viewer", "config"),
			filepath.Join(home, ".kage-viewer"),
			"kage-viewer.conf",
		}
	}

	// Load the config file[s]
	for _, cfgfile := range configfiles {
		if path, err := os.Stat(cfgfile); !os.IsNotExist(err) {
			if !path.IsDir() {
				if err := kloader.Load(file.Provider(cfgfile), toml.Parser()); err != nil {
					return nil, fmt.Errorf("error loading config file: %w", err)
				}
			}
		} // else: we ignore the file if it doesn't exists
	}

	// command line setup
	if err := kloader.Load(posflag.Provider(flagset, ".", kloader), nil); err != nil {
		return nil, fmt.Errorf("error loading flags: %w", err)
	}

	// fetch values
	conf := &Config{}
	if err := kloader.Unmarshal("", &conf); err != nil {
		return nil, fmt.Errorf("error unmarshalling: %w", err)
	}

	if err := SanitiyCheck(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func SanitiyCheck(conf *Config) error {
	if len(conf.Image) < 1 {
		return fmt.Errorf("at least 1 image must be specified")
	}

	if len(conf.Image) > 4 {
		return fmt.Errorf("only 4 images can be specified")
	}

	if conf.Shader == "" {
		return fmt.Errorf("shader file must be specified")
	}

	geoerr := "window geometry must be specified as NUMBERxNUMBER, e.g. 32x32"

	geo := strings.Split(conf.Geo, "x")
	if len(geo) != 2 {
		return fmt.Errorf(geoerr)
	}

	w, errw := strconv.Atoi(geo[0])
	h, errh := strconv.Atoi(geo[1])
	if errw != nil || errh != nil {
		return fmt.Errorf(geoerr)
	}

	conf.Width = w
	conf.Height = h

	poserr := "shader position must be specified as NUMBERxNUMBER, e.g. 32x32"

	pos := strings.Split(conf.Posision, "x")
	if len(geo) != 2 {
		return fmt.Errorf(poserr)
	}

	x, errx := strconv.Atoi(pos[0])
	y, erry := strconv.Atoi(pos[1])
	if errx != nil || erry != nil {
		return fmt.Errorf(poserr)
	}

	conf.X = x
	conf.Y = y

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
