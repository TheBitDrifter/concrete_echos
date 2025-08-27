package components

type SimpleNotification struct {
	StartedTick                   int
	Title, Body                   string
	DisplayedTitle, DisplayedBody string
	RevealStarted                 int
	IsFinished                    bool
	IsTitleFinished               bool
	PaddingX, PaddingY            float64
	TitleMaxWidth                 float64
	BodyMaxWidth                  float64
	Wrapped                       bool
}
