package gamelogic

import "log/slog"

// 判定牌的颜色，
// 红 1
// 黄 2
// 蓝 3
// 绿 4
// 功能牌 5
func JustifyCardColor(card int) int {
	color := card / 100
	slog.Debug("Card Color Justify", "card", card, "color", color)
	return color
}

// 是否为万能牌（Wild / +4）
func IsWildCard(card int) bool {
	return card/100 == 5
}

// 是否为 +4 王牌
func IsPlusFour(card int) bool {
	return card >= 511 && card <= 514
}

// 是否为禁止牌（Skip）
func IsSkipCard(card int) bool {
	if IsWildCard(card) {
		return false
	}
	lastTwo := card % 100
	return lastTwo == 21 || lastTwo == 31
}

// 是否为反转牌（Reverse）
func IsReverseCard(card int) bool {
	if IsWildCard(card) {
		return false
	}
	lastTwo := card % 100
	return lastTwo == 22 || lastTwo == 32
}

// 是否为 +2 牌（Draw 2）
func IsDrawTwoCard(card int) bool {
	if IsWildCard(card) {
		return false
	}
	lastTwo := card % 100
	return lastTwo == 23 || lastTwo == 33
}

// 是否为数字牌（0-9）
func isNumberCard(card int) bool {
	if IsWildCard(card) || IsSkipCard(card) || IsReverseCard(card) || IsDrawTwoCard(card) {
		return false
	}
	return card >= 100 && card <= 499
}

// 数字牌的点数，0-9
func cardNumber(card int) int {
	return (card % 100) % 10
}
