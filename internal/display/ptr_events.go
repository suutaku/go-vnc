package display

import (
	"github.com/go-vgo/robotgo"
	"github.com/suutaku/go-vnc/internal/types"
)

func (d *Display) servePointerEvent(ev *types.PointerEvent) {
	robotgo.Move(int(ev.X), int(ev.Y))
	btns := make(map[string]bool)
	for mask, maskType := range btnMasks {
		btns[maskType] = nthBitOf(ev.ButtonMask, mask) == 1
		if btns[maskType] {
			robotgo.MouseDown(robotGoKeyNames[maskType])
		} else {
			robotgo.MouseUp(robotGoKeyNames[maskType])
		}
	}

	// for k, v := range btns {
	// 	switch k {
	// 	case "left", "middle", "right", "scroll-up", "scroll-down", "scroll-left", "scroll-right":
	// 		if v {
	// 			robotgo.MouseDown(robotGoKeyNames[k])
	// 		} else {
	// 			robotgo.MouseUp(robotGoKeyNames[k])
	// 		}
	// 	case "unhandled":
	// 	}
	// }
}

var btnMasks = map[int]string{
	0: "left",
	1: "middle",
	2: "right",
	3: "scroll-up",
	4: "scroll-down",
	5: "scroll-left",
	6: "scroll-right",
	7: "unhandled",
}

var robotGoKeyNames = map[string]string{
	"left":         "left",
	"middle":       "center",
	"right":        "right",
	"scroll-up":    "wheelUp",
	"scroll-down":  "wheelDown",
	"scroll-left":  "wheelLeft",
	"scroll-right": "wheelRight",
	"unhandled":    "unhandled",
}

func nthBitOf(bit uint8, n int) uint8 {
	return (bit & (1 << n)) >> n
}
