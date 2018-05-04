# Notes

You can decode the hex that's output in the EDID section from `xrandr --props`:

```
$ echo "00ffffffffffff004d1053140000000028190104a52313780ede50a3544c99260f5054000000010101010101010101010101010101011a3680a070381f40302035005ac210000018000000000000000000000000000000000000000000fe00313230334d814c513135364d31000000000002410328001200000a010a2020004a" | xxd -r -p | parse-edid
```

Vendor names are available in `/usr/share/hwdata/pnp.ids`. This file is 
available in the Ubuntu release of the `hwdata` package, but not Arch's.

You can also get the EDID from a file. Each monitor should have it's own 
folder under:

```
$ cat /sys/devices/pci0000:00/0000:00:02.0/drm/card0/card0-<MON>/edid`
```

This can then be parsed with `parse-edid`.

Of course, for identify a monitor as uniquely as possible, just using an MD5
sum of the EDID will probably be enough. This can be checked on a device with
2 displays of the same kind. Even then, it won't matter if it can't 
distinguish between those 2, i3adc will just need to remember which port either
was in and act the same for that port and display type.

Looks like you can also get the edid from a shorter path:

```
$ cat /sys/class/drm/card0-<MON>/edid
```

From there, we could just make a hash of the EDID, we don't need to parse it at all. With that, and 
the port we'd know all we need to know about a display, and could print debug info for getting more
information about a display using `read-edid`. Gotta keep in mind too, this is at least at first 
going to be a fairly specialised tool...

## Events

* Which changes trigger events in i3? We need to know when a new primary monitor is selected, or 
when resolutions or positions change, or when something is unplugged or plugged in, etc. If a 
monitor is not in use, do these events still happen, or is that outside of i3, i.e. in X?
