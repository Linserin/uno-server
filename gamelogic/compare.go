package gamelogic

import "log/slog"

// compare 判断当前牌是否可以接在上一张牌之后。
// activeColor 为当前生效颜色（万能牌出牌后由出牌者指定，保存在 ActionData.Color 中）。
func compare(activeColor int, thisCard int, lastCard int) bool {
	if thisCard == 1000 {
		slog.Info("Player skip, next turn")
		return true
	}
	// 万能牌永远可以出
	if IsWildCard(thisCard) {
		slog.Info("Wild card is always playable")
		return true
	}
	// 匹配当前生效颜色
	if JustifyCardColor(thisCard) == activeColor {
		slog.Info("颜色匹配，返回 True")
		return true
	}
	// 数字牌匹配点数
	if isNumberCard(thisCard) && isNumberCard(lastCard) && cardNumber(thisCard) == cardNumber(lastCard) {
		slog.Info("数字匹配，返回 True")
		return true
	}
	// 功能牌匹配同类功能
	if (IsSkipCard(thisCard) && IsSkipCard(lastCard)) ||
		(IsReverseCard(thisCard) && IsReverseCard(lastCard)) ||
		(IsDrawTwoCard(thisCard) && IsDrawTwoCard(lastCard)) {
		slog.Info("功能匹配，返回 True")
		return true
	}
	slog.Info("颜色、数字与功能均不匹配")
	return false
}
