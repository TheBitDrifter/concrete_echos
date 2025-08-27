package animations

import "github.com/TheBitDrifter/bappa/blueprint/client"

var DefeatScreenAnim = client.AnimationData{
	Name:        "defeatScreen",
	RowIndex:    0,
	FrameCount:  20,
	Speed:       10,
	FrameWidth:  640,
	FrameHeight: 360,
	Freeze:      true,
}

var HomeScreenAnim = client.AnimationData{
	Name:        "homeScreen",
	RowIndex:    0,
	FrameCount:  15,
	Speed:       10,
	FrameWidth:  640,
	FrameHeight: 360,
}
