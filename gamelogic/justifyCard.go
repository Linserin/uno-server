package gamelogic

import "log/slog"

// 判定牌的颜色，
// 红 1
// 黄 2
// 蓝 3
// 绿 4
// 功能牌 5
func JustifyCardColor(card int) int {
	var color int
	for card > 100 {
		card -= 100
		color++
	}
	slog.Debug("Card Color Justify", "card", card, "color", color)
	return color
}
