package gamelogic

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"strconv"
)

type GameActionData struct {
	Ban       bool `json:"ban"`       // ban 为禁止（禁止下一个出牌）
	PlusTwo   int  `json:"plusTwo"`   // plusTwo 为 +2, 后方要不加二要不叠加给下一个
	Color     int  `json:"color"`     // 对使用 Wide 时进行记录
	ClockWise bool `json:"clockWise"` // 方向， true 为顺时针 否则为逆时针
	Turn      int  `json:"turn"`      // 轮到谁来
	LastCard  int  `json:"lastCard"`  // 上一张卡
	PlusFour  int  `json:"plusFour"`  // 王牌叠加 4
	Players   int  `json:"players"`   // 游戏人数（2-4）
}

// 从切片里删除第一个等于 card 的元素，返回新切片和是否找到
func removeCardFromHand(cards []int, card int) ([]int, bool) {
	for i, c := range cards {
		if c == card {
			// append 会就地覆盖，返回的新切片长度 -1
			return append(cards[:i], cards[i+1:]...), true
		}
	}
	return cards, false
}

// 判断牌是否合法拥有
func whetherCardsValid_Have(roomID int, player int, card int) error {
	// 如果 card 是 1000，表示跳过出牌，无需拥有
	if card == 1000 {
		return nil
	}
	allCards := []int{
		// Red 100：25
		100,
		101, 102, 103, 104, 105, 106, 107, 108, 109,
		111, 112, 113, 114, 115, 116, 117, 118, 119,
		121, 131, // Skip
		122, 132, // Reverse
		123, 133, // Draw 2

		// Yellow 200：25
		200,
		201, 202, 203, 204, 205, 206, 207, 208, 209,
		211, 212, 213, 214, 215, 216, 217, 218, 219,
		221, 231, // Skip
		222, 232, // Reverse
		223, 233, // Draw 2

		// Blue 300：25
		300,
		301, 302, 303, 304, 305, 306, 307, 308, 309,
		311, 312, 313, 314, 315, 316, 317, 318, 319,
		321, 331, // Skip
		322, 332, // Reverse
		323, 333, // Draw 2

		// Green 400：25
		400,
		401, 402, 403, 404, 405, 406, 407, 408, 409,
		411, 412, 413, 414, 415, 416, 417, 418, 419,
		421, 431, // Skip
		422, 432, // Reverse
		423, 433, // Draw 2

		// Other Cards
		511, 512, 513, 514, // +4 & Wild
		521, 522, 523, 524, // Wild Change Color
	}
	isLegally := false
	for _, cardInRepo := range allCards {
		if card == cardInRepo {
			isLegally = true
			break
		}
	}
	if !isLegally {
		return errors.New("Illegal card ID")
	}
	// 阅读文件并且解析
	room, err := os.Open(strconv.Itoa(roomID) + ".json")
	if err != nil {
		slog.Error("[gamelogic/main.go] An error occurred", "error", err)
		return err
	}
	defer room.Close()
	data, err := io.ReadAll(room)
	if err != nil {
		slog.Error("[gamelogic/main.go] An error occurred", "error", err)
		return err
	}
	var gameData UNO_Cards
	err = json.Unmarshal(data, &gameData)
	if err != nil {
		slog.Error("[gamelogic/main.go] An error occurred", "error", err)
		return err
	}
	// 判定玩家手中是否存在对应卡牌
	found := false
	switch player {
	case 1:
		for _, cardsInPlayer := range gameData.First {
			if card == cardsInPlayer {
				found = true
			}
		}
	case 2:
		for _, cardsInPlayer := range gameData.Sec {
			if card == cardsInPlayer {
				found = true
			}
		}
	case 3:
		for _, cardsInPlayer := range gameData.Third {
			if card == cardsInPlayer {
				found = true
			}
		}
	case 4:
		for _, cardsInPlayer := range gameData.Fourth {
			if card == cardsInPlayer {
				found = true
			}
		}
	}
	// found bool 为是否存在
	if found {
		slog.Info("Found Card in player", "roomID", roomID, "player", player, "card", card)
		return nil
	}
	slog.Info("Card Not Found in player", "roomID", roomID, "player", player, "card", card)
	return errors.New("Card Not Found in player")
}

// playerCount 返回游戏人数，旧存档没有 players 字段时按 4 人处理
func playerCount(actionData GameActionData) int {
	if actionData.Players >= 2 && actionData.Players <= 4 {
		return actionData.Players
	}
	return 4
}

// advanceTurn 按当前方向前进 steps 个玩家
func advanceTurn(actionData *GameActionData, steps int) {
	count := playerCount(*actionData)
	turn := actionData.Turn
	for i := 0; i < steps; i++ {
		if actionData.ClockWise {
			turn++
			if turn > count {
				turn = 1
			}
		} else {
			turn--
			if turn < 1 {
				turn = count
			}
		}
	}
	actionData.Turn = turn
}

// drawCards 从牌库给指定玩家摸 count 张牌
func drawCards(gameData *UNO_Cards, player int, count int) {
	for i := 0; i < count && len(gameData.Repo) > 0; i++ {
		card := gameData.Repo[0]
		gameData.Repo = gameData.Repo[1:]
		switch player {
		case 1:
			gameData.First = append(gameData.First, card)
		case 2:
			gameData.Sec = append(gameData.Sec, card)
		case 3:
			gameData.Third = append(gameData.Third, card)
		case 4:
			gameData.Fourth = append(gameData.Fourth, card)
		}
	}
}

// 移除玩家手中的牌，返回是否成功
func removeCardFromPlayer(gameData *UNO_Cards, player int, card int) bool {
	var removed bool
	switch player {
	case 1:
		gameData.First, removed = removeCardFromHand(gameData.First, card)
	case 2:
		gameData.Sec, removed = removeCardFromHand(gameData.Sec, card)
	case 3:
		gameData.Third, removed = removeCardFromHand(gameData.Third, card)
	case 4:
		gameData.Fourth, removed = removeCardFromHand(gameData.Fourth, card)
	}
	return removed
}

func Logic_Action_main(roomID int, player int, card int, color int) error {
	slog.Debug("[gamelogic/main.go] Card Action.", "roomID", roomID, "player", player, "card", card, "color", color)
	if card == 0 {
		return errors.New("card cannot be 0")
	}
	actionData, err := ReadGameActionData(roomID)
	if err != nil {
		slog.Error("[gamelogic/main.go] An error occurred", "error", err)
		return err
	}
	if actionData.Turn != player {
		slog.Error("[gamelogic/main.go] An error occurred", "error", errors.New("Not your turn"))
		return errors.New("Not your turn")
	}

	// ID 1000 表示跳过（Pass）：不改变 LastCard，只移动回合。
	// 若当前存在 +2 / +4 罚牌，则摸走对应的牌并结束自己的回合。
	if card == 1000 {
		slog.Info("Player skip, next turn", "roomID", roomID, "player", player)
		if actionData.PlusTwo > 0 || actionData.PlusFour > 0 {
			gameData, err := ReadGameData(roomID)
			if err != nil {
				slog.Error("[gamelogic/main.go] An error occurred", "error", err)
				return err
			}
			drawCards(&gameData, player, actionData.PlusTwo+actionData.PlusFour)
			actionData.PlusTwo = 0
			actionData.PlusFour = 0
			if err := WriteGameData(roomID, gameData); err != nil {
				slog.Error("[gamelogic/main.go] An error occurred", "error", err)
				return err
			}
		}
		advanceTurn(&actionData, 1)
		return WriteGameActionData(roomID, actionData)
	}

	err = whetherCardsValid_Have(roomID, player, card)
	if err != nil {
		slog.Error("[gamelogic/main.go] An error occurred", "error", err)
		return err
	}

	// 有罚牌待结算时，只能叠加同类罚牌，否则需要用 1000 摸牌跳过
	if actionData.PlusTwo > 0 && !(IsDrawTwoCard(card) && JustifyCardColor(card) == actionData.Color) {
		return errors.New("Must stack +2 or pass")
	}
	if actionData.PlusFour > 0 && !IsPlusFour(card) {
		return errors.New("Must stack +4 or pass")
	}

	if !compare(actionData.Color, card, actionData.LastCard) {
		slog.Error("[gamelogic/main.go] An error occurred", "error", errors.New("Card Not Legal"))
		return errors.New("Card Not Legal")
	}

	gameData, err := ReadGameData(roomID)
	if err != nil {
		slog.Error("[gamelogic/main.go] An error occurred", "error", err)
		return err
	}

	if !removeCardFromPlayer(&gameData, player, card) {
		slog.Error("[gamelogic/main.go] An error occurred", "error", errors.New("Card Not Found in player"))
		return errors.New("Card Not Found in player")
	}

	// 更新当前生效颜色
	if IsWildCard(card) {
		if color >= 1 && color <= 4 {
			actionData.Color = color
		} else if actionData.Color < 1 || actionData.Color > 4 {
			actionData.Color = 1
		}
	} else {
		actionData.Color = JustifyCardColor(card)
	}

	// 不同功能牌的逻辑
	switch {
	case IsSkipCard(card):
		// 禁止：下一位玩家暂停出牌，直接轮到再下家
		advanceTurn(&actionData, 2)
	case IsReverseCard(card):
		// 反转：改变方向；两人对战时等同于禁止
		if playerCount(actionData) == 2 {
			advanceTurn(&actionData, 2)
		} else {
			actionData.ClockWise = !actionData.ClockWise
			advanceTurn(&actionData, 1)
		}
	case IsDrawTwoCard(card):
		// +2：累加罚牌，轮到下家决定叠加还是摸牌
		actionData.PlusTwo += 2
		advanceTurn(&actionData, 1)
	case IsPlusFour(card):
		// +4：累加罚牌，轮到下家决定叠加还是摸牌
		actionData.PlusFour += 4
		advanceTurn(&actionData, 1)
	default:
		advanceTurn(&actionData, 1)
	}

	// 更新 LastCard
	actionData.LastCard = card
	// 写入数据
	err = WriteGameData(roomID, gameData)
	if err != nil {
		slog.Error("[gamelogic/main.go] An error occurred", "error", err)
		return err
	}
	err = WriteGameActionData(roomID, actionData)
	if err != nil {
		slog.Error("[gamelogic/main.go] An error occurred", "error", err)
		return err
	}
	return nil
}
