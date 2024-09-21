# AncientPlotter

This repository is a work in progress to revive a COGI CT-630 (and all its OEM-Name-Clones) serial Plotter.

I've found some true and untrue informations as well as some Python-software in the web, however the found software always enforced me to painfully deal with the Python dependencies on MacOS. So I assembled what needed to be done and wrote this little tool. Thanks to the past [Python developers to discover the CT-630 endsequence-codes](https://gitli.stratum0.org/chrissi/cogidraw/-/blob/master/cogidraw.py?ref_type=heads#L58) and publishing them on the web.

## Additional motivation

[Inkscape](https://inkscape.org) in the version 1.3.2 at time of writing can save to `.hpgl` files (this is my way to convert a svg to a hpgl). However these have proven to be incompatible with my version of the CT-630. Therefore `ancientPlotter` does the following to succeed:
- parse the given hpgl-file and split out all incompatible series to single instructions
- send one instruction at once (there is reportedly a buffer problem with this plotter)
- send the end-sequence

## Contribution

Feel free to file a PR if you like. I will try to process ASAP.

## Binary distribution

This repository has binaries for Windows, Linux and MacOS each compiled for arm64 (intel) or arm64 (Apple / Raspberry) in the [bin](bin) folder. However, at least MacOS won't allow you to execute a binary directly loaded from the internet - which is pretty understandable.

## Using docker - WIP

I will - somewhere in the future - add a github-CI soon which will provide a docker image for the above mentioned architectures.

## Building on your own

The safest way is to build the executable on your own. Here is how on Linux and MacOS. I dont have a clue how to do anything on Windows. If you found out, feel free to open a PR to this readme with instructions 😚

- [Install go](https://go.dev/doc/install) in your favourite way.
- Install [git](https://git-scm.com/) if not already done. 
- Checkout this repository and build the binary.

```bash
# go and git are already installed
git clone https://github.com/Dirk007/ancientPlotter.git
cd ancientPlotter
go build .
```

go will hopefully download all dependencies and place you a nice `ancientPlotter` executable in the current directory. You should be able to run it straight ahead

```bash
./ancientPlotter
```

### Building with `dagger`

Invoke 
```bash
dagger call build --src=. 
```

## Usage

`ancientPlotter [--serial-device <device>] [--dry-run] [--print-only] [--serve] [--port <port> (default 11175)] filename.hpgl`

`--serve`

Spin up this tool as a backend serving the required HTML and Websockets to drive your plotter.

You can reach the service by invoking http://127.0.0.1:11175/ (or whatever port you assigned). The documents will be served from [assets](./assets/). You can modify them at will *at runtime*. They will be freshly read and rendered on each onvocation.

Please don't blame me for the bad quality - I am everything but a frontend engineer.

`--port`

If you specifiy this parameter and also use the `--serve` option, the default port `11175` will be overriden with this value.

`--filename`

The only required parameter if you use this tool as a CLI, specifies which file you want to plot. The path can be absolute or relative - it does not matter. If you `serve` this parameter will simply be ignored.

Ony the first given file will be plotted.

`--serial-device`

Specify the path to your serial-device where the Plotter is connected to. For example `--serial-device /dev/tty.usbserial-10`. The name of the port highly depends on your system but something with `usb` should point you to the right direction.

This parameter is **optional** - if you don't specify the serial-device, `ancientPlotter` tries to *guess* (relatively dumb) the right port. This works for me on my Mac, but I have no idea how stable it is on other machines. Maybe you can leave some feedback about that.

`--dry-run`

If you specify this parameter, the instructions will be pached on the fly to **not** put the pen (cutter) down but always keep it hovering. I made this to not ruin too much foil while exploring the plotters functionality / scale. Maybe it is useful for someone else so I kept it. You may use it to see if your scaling is more or less correct. **WARNING** I have no idea if this will cause damage to the plotter if done too often. I am relatively sure that the plotter is designed to work with the resistence of the foil. But not sure. However, you've been warned.

`--print-only`

If you specify this parameter, the instructions will **not** be sent to the Plotter but printed to the console instead. It is just for educational reasons and for debugging.

## Other resources in the web

You make take a look at [stratum0.org](https://stratum0.org/wiki/Cogi_CT-630) which initially brought me some knowledge about this little diva of a Plotter.

