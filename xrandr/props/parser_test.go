package props

import "testing"

func BenchmarkParser_ParseProps(b *testing.B) {
	parser := NewParser(false)

	for i := 0; i < b.N; i++ {
		_, err := parser.ParseProps(propsOutput)
		if err != nil {
			b.Error(err)
		}
	}
}

var propsOutput = []byte(`Screen 0: minimum 320 x 200, current 1920 x 1080, maximum 8192 x 8192
eDP-1 connected primary 1920x1080+0+0 (normal left inverted right x axis y axis) 346mm x 194mm
	EDID: 
		00ffffffffffff004d10531400000000
		28190104a52313780ede50a3544c9926
		0f505400000001010101010101010101
		0101010101011a3680a070381f403020
		35005ac2100000180000000000000000
		00000000000000000000000000fe0031
		3230334d814c513135364d3100000000
		0002410328001200000a010a2020004a
	scaling mode: Full aspect 
		supported: Full, Center, Full aspect
	Broadcast RGB: Automatic 
		supported: Automatic, Full, Limited 16:235
	link-status: Good 
		supported: Good, Bad
	CONNECTOR_ID: 59 
		supported: 59
	non-desktop: 0 
		range: (0, 1)
   1920x1080     59.93*+
   1680x1050     59.95    59.88  
   1400x1050     59.98  
   1600x900      59.99    59.94    59.95    59.82  
   1280x1024     60.02  
   1400x900      59.96    59.88  
   1280x960      60.00  
   1440x810      60.00    59.97  
   1368x768      59.88    59.85  
   1280x800      59.99    59.97    59.81    59.91  
   1280x720      60.00    59.99    59.86    59.74  
   1024x768      60.04    60.00  
   960x720       60.00  
   928x696       60.05  
   896x672       60.01  
   1024x576      59.95    59.96    59.90    59.82  
   960x600       59.93    60.00  
   960x540       59.96    59.99    59.63    59.82  
   800x600       60.00    60.32    56.25  
   840x525       60.01    59.88  
   864x486       59.92    59.57  
   700x525       59.98  
   800x450       59.95    59.82  
   640x512       60.02  
   700x450       59.96    59.88  
   640x480       60.00    59.94  
   720x405       59.51    58.99  
   684x384       59.88    59.85  
   640x400       59.88    59.98  
   640x360       59.86    59.83    59.84    59.32  
   512x384       60.00  
   512x288       60.00    59.92  
   480x270       59.63    59.82  
   400x300       60.32    56.34  
   432x243       59.92    59.57  
   320x240       60.05  
   360x202       59.51    59.13  
   320x180       59.84    59.32  
DP-1 disconnected (normal left inverted right x axis y axis)
	Broadcast RGB: Automatic 
		supported: Automatic, Full, Limited 16:235
	audio: auto 
		supported: force-dvi, off, auto, on
	link-status: Good 
		supported: Good, Bad
	CONNECTOR_ID: 66 
		supported: 66
	non-desktop: 0 
		range: (0, 1)
HDMI-1 disconnected (normal left inverted right x axis y axis)
	aspect ratio: Automatic 
		supported: Automatic, 4:3, 16:9
	Broadcast RGB: Automatic 
		supported: Automatic, Full, Limited 16:235
	audio: auto 
		supported: force-dvi, off, auto, on
	link-status: Good 
		supported: Good, Bad
	CONNECTOR_ID: 71 
		supported: 71
	non-desktop: 0 
		range: (0, 1)
DP-2 disconnected (normal left inverted right x axis y axis)
	Broadcast RGB: Automatic 
		supported: Automatic, Full, Limited 16:235
	audio: auto 
		supported: force-dvi, off, auto, on
	link-status: Good 
		supported: Good, Bad
	CONNECTOR_ID: 74 
		supported: 74
	non-desktop: 0 
		range: (0, 1)
HDMI-2 disconnected (normal left inverted right x axis y axis)
	aspect ratio: Automatic 
		supported: Automatic, 4:3, 16:9
	Broadcast RGB: Automatic 
		supported: Automatic, Full, Limited 16:235
	audio: auto 
		supported: force-dvi, off, auto, on
	link-status: Good 
		supported: Good, Bad
	CONNECTOR_ID: 78 
		supported: 78
	non-desktop: 0 
		range: (0, 1)
`)
