# i3adc

Automatic display configuration for i3 window manager, written in Go.

## Usage

There isn't a pre-built set of binaries just yet, but as long as you have Go installed, you can 
build it yourself (you'll need to [install dep][1]).

```
$ go get -u -v github.com/seeruk/i3adc
$ cd $GOPATH/src/github.com/seeruk/i3adc
$ dep ensure -v
$ go install -v ./cmd/...
```

Once you have it installed, if you've got your `$GOPATH` set up properly, and `$GOPATH/bin` on your
`$PATH`, you should just be able to use the `i3adc` command.

```
$ i3adc 
```

Personally, I have a line in my `.xinitrc` to run `i3adc` when I log in. Use whatever works for you.

## Events

i3adc receives display events from i3's IPC. They don't usually come attached with any information.
Display events will occur when displays are plugged in or unplugged, and when display configuration 
is changed (e.g. manually via `xrandr`, or maybe via some kind of graphical application like 
ARandR).

When an event occurs, it can be handled in one of 3 different ways:

1. If a new set of displays is detected (i.e. the combination of all connected displays is new); 
then all connected displays will be enabled, and set to their preferred mode. Their positions will 
also be reset. Positioning is based off of the order that the displays are sent from X, and each 
display will be to the right of the previous display (in one long row). This layout will then be 
saved.

2. If the set of connected displays doesn't change, but some other settings change (e.g. modes, 
positions, rotation, reflection, so on...), then i3adc will update the saved configuration for that 
layout. This configuration is stored in the file `~/.i3adc/i3adc.db`. 

3. If a display is plugged in, or unplugged, and the new set of displays has a configuration saved 
already; then i3adc will apply the same configuration that was used last time those displays were 
connected.

As an example, let's say you have a laptop, and 2 external monitors. You're running i3, and have 
only just plugged those displays in - so they're not active right now, or even enabled (i.e. they're
on standby). If you started i3adc for the first time, all 3 displays would turn on, they would be 
set to their preferred mode, and positioned so that none of them overlap.

If you then decided you don't want to use the laptop display when you have your 2 external displays 
connected, you could disable the laptop display by turning it off with `xrandr`. i3adc would notice
this change and update the saved layout. If you then disconnected both of your external displays,
your laptop display would be re-enabled. If you then plugged both of your external displays in again
the laptop display would turn off, and the two external displays would turn on and return to the 
same configuration they had last time they were connected.

## License 

MIT


[1]: https://golang.github.io/dep/docs/installation.html
