# Notes

## Display Configurations

Events that are relevant to reading, creating, or updating i3adc state are as follows:

* When i3adc starts.
* When i3adc detects a new output has been connected.
* When i3adc detects an existing output has been disconnected.

### On Output Event

When an event occurs, i3adc will need to react in only 2 different ways; read existing 
configuration, or create/update new configuration.

#### Read / Create

When i3adc is starting up, when a new display is connected, or when a display is disconnected, i3adc
should look at the currently connected outputs, create their unique hash, then look up the stored 
configuration for that hash. 

If there is no configuration stored for the currently connected outputs, i3adc should activate the
preferred mode for each output, and then save that as the initial configuration. This behaviour is
fairly primitive for the time being to keep the "MVP" small but functional.

If there is existing configuration, that should be applied, if it is possible for it to be applied 
after some validation.

* What about primary display? One must always be primary, does xrandr handle that for us and just
make sure that one is always the primary display? If not, we can just pick the first. This can 
always be changed later on.
* What about offsets? Displays that overlay behave... interestingly. If a new configuration is made,
it's easy to ensure that we'll offset them in the right order. But then, we've already seen this 
happen before with other operating systems. It's not a big deal for initial configuration.

#### Updated

When display configuration is changed, but the connected devices have not changed, we should be able
to detect this, and then update configuration. This should never occur on startup. When i3adc starts
it should record the currently active "profile". If that profile doesn't change when an event 
occurs, it should record the new configuration and save that instead. 

## Output Identification

Each output should be uniquely identifiable, reliably. One of the main properties of an output that
can be used to achieve this is the EDID. Realistically though, all properties could be hashed 
together to identify an output, which would then help in the case that an output has no EDID for 
some reason.

From the hashes of each individual output, a hash can easily be created for all connected outputs. 
This can then be used to; for example, store that configuration under a key that should be 
reproducible.

## Hooks

When any change occurs in i3adc, scripts should be able to be called as "hooks". This can be used
for example to inform software like Polybar to restart, or for a new primary display to have the 
primary bar. More specifics to come...
