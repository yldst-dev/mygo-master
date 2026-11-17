package main

import (
	"time"

	"github.com/egoist/mygo/ui"
)

var (
	pressMotion = ui.Animation{Duration: 100 * time.Millisecond}
	fadeMotion  = ui.Animation{Duration: 120 * time.Millisecond}
	quickMotion = ui.Animation{Duration: 150 * time.Millisecond}
	themeMotion = ui.Animation{Duration: 200 * time.Millisecond}
	markSpring  = ui.Animation{Spring: ui.SpringOptions{Duration: 180 * time.Millisecond}}
	pillMotion  = ui.Animation{Spring: ui.SpringOptions{Duration: 240 * time.Millisecond}}

	pageEnter  = ui.ElementTransition{Colors: true, Duration: 180 * time.Millisecond, Enter: &ui.Motion{Y: 6, Opacity: 0.35}}
	markMotion = ui.ElementTransition{Position: true, Animation: &markSpring, Enter: &ui.Motion{Scale: 0.9}, Exit: &ui.Motion{Scale: 0.9}}
	popMotion  = ui.ElementTransition{Position: true, Duration: 120 * time.Millisecond, Enter: &ui.Motion{Scale: 0.6}, Exit: &ui.Motion{Scale: 0.6}}
	dotExit    = ui.ElementTransition{Colors: true, Duration: 150 * time.Millisecond, Exit: &ui.Motion{Scale: 0.4}}
	growFit    = ui.ElementTransition{Size: true, Animation: &pillMotion}
	followFit  = ui.ElementTransition{Position: true, Animation: &pillMotion}
	tipFit     = ui.ElementTransition{Size: true, Animation: &markSpring}
	rollDigits = ui.ContentTransition{Numeric: true, Animation: &pillMotion}
	tipDigits  = ui.ContentTransition{Numeric: true, Animation: &fadeMotion}
)

var bellRing = []ui.Keyframe{
	{Value: 0},
	{At: 70 * time.Millisecond, Value: 14},
	{At: 150 * time.Millisecond, Value: -10},
	{At: 230 * time.Millisecond, Value: 6},
	{At: 300 * time.Millisecond, Value: -3},
	{At: 360 * time.Millisecond, Value: 0},
}

func progress(e ui.Element, key string, on bool, motion ui.Animation) float32 {
	target := float32(0)
	if on {
		target = 1
	}
	return e.AnimateTo(key, target, motion)
}

func pressable(e ui.Element, scale float32) ui.Element {
	var pose ui.Pose
	if e.Pressed() {
		pose.Scale = scale
	}
	return e.AnimatePose(pose, pressMotion)
}
