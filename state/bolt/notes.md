# i3adc: Bolt Notes

* If migrations need to occur to make a key/value store work, then they should happen when that type
of store is used, not in general. This is one of the things that complicated the storage in Tid. It
had too many layers caused by an abstraction that was only necessary because of Bolt.
