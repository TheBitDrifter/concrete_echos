package animations

import "github.com/TheBitDrifter/bappa/blueprint/client"

var DialogueOpenAnimation = client.AnimationData{
	Name:        "dialogueOpen",
	RowIndex:    2,
	FrameCount:  4,
	Speed:       5,
	FrameWidth:  395,
	FrameHeight: 78,
	Freeze:      true,
}

var DialogueCloseAnimation = client.AnimationData{
	Name:        "dialogueClose",
	RowIndex:    1,
	FrameCount:  3,
	Speed:       5,
	FrameWidth:  395,
	FrameHeight: 78,
	Freeze:      true,
}

var DialogueNext = client.AnimationData{
	Name:        "dianext",
	RowIndex:    0,
	FrameCount:  4,
	Speed:       8,
	FrameWidth:  7,
	FrameHeight: 8,
	Freeze:      false,
}
