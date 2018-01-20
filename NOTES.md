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
