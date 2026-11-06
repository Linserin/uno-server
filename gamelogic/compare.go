package gamelogic

import "log/slog"

func compare(thisCard int, lastCard int) bool {
	// 匹配颜色
	thisColor := JustifyCardColor(thisCard)
	lastColor := JustifyCardColor(lastCard)
	if lastColor == thisColor {
		slog.Info("颜色匹配，返回 True")
		return true
	}
	slog.Info("颜色不匹配，下一步...")
	// 匹配个位数字是否相等
	thisNum := thisCard % 10
	lastNum := lastCard % 10
	if thisNum == lastNum {
		slog.Info("数字匹配，下一步...")
		return true
	}
	slog.Info("数字不匹配，下一步...")
	// 王牌与 Wild
	if thisCard%100 == 5 {
		return true
	}
	return false
}
