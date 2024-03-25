# kage-viewer - Viewer for shaders written in Kage, similar to glslviewer

![Logo](https://github.com/TLINDEN/kageviewer/blob/main/.github/assets/logo.png)

[![License](https://img.shields.io/badge/license-GPL-blue.svg)](https://github.com/tlinden/kage-viewer/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tlinden/kage-viewer)](https://goreportcard.com/report/github.com/tlinden/kage-viewer) 

This   little  tool   can  be   used  to   test  shaders   written  in
[Kage](https://ebitengine.org/en/documents/shader.html), a shader meta
language for
[Ebitengine](https://github.com/hajimehoshi/ebiten). kage-viewer
reloads changed assets, which allows you to develop your shader and
see live, how it responds to your changes. If loading fails, an error
will be printed to STDOUT. The same applies for images.

## Screenshot

![Screenshot](https://github.com/TLINDEN/kageviewer/blob/main/.github/assets/screenshot.png)

## Installation

The tool does not have any dependencies.  Just download the binary for
your platform from the releases page and you're good to go.

### Installation using a pre-compiled binary

Go to the [latest release page](https://github.com/TLINDEN/kage-viewer/releases/latest)
and look for your OS and platform. There are two options to install the binary:

Directly     download     the     binary    for     your     platform,
e.g. `kage-viewer-linux-amd64-0.0.2`, rename it to `kage-viewer` (or whatever
you like more!)  and put it into  your bin dir (e.g. `$HOME/bin` or as
root to `/usr/local/bin`).

Be sure  to verify  the signature  of the binary  file. For  this also
download the matching `kage-viewer-linux-amd64-0.0.2.sha256` file and:

```shell
cat kage-viewer-linux-amd64-0.0.2.sha25 && sha256sum kage-viewer-linux-amd64-0.0.2
```
You should see the same SHA256 hash.

You  may  also download  a  binary  tarball  for your  platform,  e.g.
`kage-viewer-linux-amd64-0.0.2.tar.gz`,  unpack and  install it.  GNU Make  is
required for this:
   
```shell
tar xvfz kage-viewer-linux-amd64-0.0.2.tar.gz
cd kage-viewer-linux-amd64-0.0.2
sudo make install
```

### Installation from source

You will need the Golang toolchain  in order to build from source. GNU
Make will also help but is not strictly neccessary.

If you want to compile the tool yourself, use `git clone` to clone the
repository.   Then   execute   `go    mod   tidy`   to   install   all
dependencies. Then  just enter `go  build` or -  if you have  GNU Make
installed - `make`.

To install after building either copy the binary or execute `sudo make
install`. 

# Usage

```shell
kage-viewer -h
This is kage-viewer, a shader viewer.

Usage: kage-viewer [-vd] [-c <config file>] [-g geom] [-p geom] \
       -i <image0.png> -i <image1.png> -s <shader.kage>

Options:
-c --config     <toml file>     Config file to use (optional)
-i --image      <png file>      Image to load (multiple times allowed, up to 4)
-s --shader     <kage file>     Shader to run
-g --geometry   <WIDTHxHEIGHT>  Window size
-p --position   <XxY>           Position of image0
   --map-flag   <name>          Map Flag uniform to <name>
   --map-ticks  <name>          Map Flag uniform to <name>
   --map-slider <name>          Map Flag uniform to <name>
   --map-mouse  <name>          Map Flag uniform to <name>
-d --debug                      Show debugging output
-v --version                    Show program version
```

Example usage using the provided example:

```shell
kage-viewer -g 32x32 -i example/wall.png -i example/damage.png  -s example/destruct.kg
```

Hit `SPACE` or press the left mouse button to toggle the damage
mask. Press the `UP` or `DOWN` key to adjust the damage scale.

# Uniforms

Since this is a generic viewer, you cannot (yet) use custom
uniforms. If you need this, just edit the source accordingly.

Uniforms supported so far:

- `var Flag int`: a flag which toggles between 0 and 1 by pressing
  `SPACE` or pusing the left mouse button
- `var Slider float`: a normalized float value, you can increment it
  with `UP` or `DOWN`
- `var Ticks float`: the time the game runs (ticks, not seconds!)
- `var Mouse vec2`: the current mouse position

If you want to test an existing shader and don't want to rename the
uniforms, you can map the ones provided by **kage-viewer** to custom
names using the `--map-*` options. For example:

```shell
kage-viewer -g 640x480 --map-ticks Time --map-mouse Cursor examples/shader/default.go
```

This executes the example shader in the ebitenging source repository.

# Config File

You can use a config file to store your own codes, once you found one
you like. A configfile is searched in these locations in this order:

* `/etc/kage-viewer.conf`
* `/usr/local/etc/kage-viewer.conf`
* `$HOME/.config/kage-viewer/config`
* `$HOME/.kage-viewer`

You may also specify a config file on the commandline using the `-c`
flag.

Config files are expected to be in the [TOML format](https://toml.io/en/).

Possible parameters equal the long command line options.

# TODO

- [X] Implement loading of images and shader files
- [X] Implement basic shader rendering and user input
- [ ] Add custom uniforms (maybe using lua code?)
- [x] Provide a way to respond live to shader code changes (use lua as
  well?)

# Report bugs

[Please open an issue](https://github.com/TLINDEN/kage-viewer/issues). Thanks!

# License

This work is licensed under the terms of the General Public Licens
version 3.

# Author

Copyleft (c) 2024 Thomas von Dein
