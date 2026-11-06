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
	var isLegally bool
	for _, cardInRepo := range allCards {
		if card == cardInRepo {
			isLegally = true
		}

	}
	if isLegally == false {
		return errors.New("Illegal card ID")
	}
	// 阅读文件并且解析
	room, err := os.Open(strconv.Itoa(roomID) + ".json")
	if err != nil {
		slog.Error("[gamelogic/main.go] An error occurred", "error", err)
		return err
	}
	data, err := io.ReadAll(room)
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
	} else {
		slog.Info("Card Not Found in player", "roomID", roomID, "player", player, "card", card)
		return errors.New("Card Not Found in player")
	}
}

func Logic_Action_main(roomID int, player int, card int) error {
	slog.Debug("[gamelogic/main.go] Card Action.", "roomID", roomID, "player", player, "card", card)
	if card == 0 {
		return errors.New("card cannot be 0")
	}
	err := whetherCardsValid_Have(roomID, player, card)
	if err != nil {
		slog.Error("[gamelogic/main.go] An error occurred", "error", err)
		return err
	}
	err = nil
	actionData, err := ReadGameActionData(roomID)
	if err != nil {
		slog.Error("[gamelogic/main.go] An error occurred", "error", err)
		return err
	}
	slog.Info("开始进行判断")
	isCardLegal := compare(card, actionData.LastCard)
	slog.Info("Card", "isLegal", isCardLegal)
	if !isCardLegal {
		slog.Error("[gamelogic/main.go] An error occurred", "error", errors.New("Card Not Legal"))
		return errors.New("Card Not Legal")
	}
	gameData, err := ReadGameData(roomID)
	if err != nil {
		slog.Error("[gamelogic/main.go] An error occurred", "error", err)
		return err
	}
	gameActionData, err := ReadGameActionData(roomID)
	// 加一轮
	gameActionData.Turn += 1
	if gameActionData.Turn == 5 {
		gameActionData.Turn = 1
	}
	// 如果是 Skip 直接加一轮
	if card%100 == 21 || card%100 == 31 {
		gameActionData.Turn += 1
	}
	var done bool
	// 移除玩家的对应牌
	switch player {
	case 1:
		for _, cardsInPlayer := range gameData.First {
			if card == cardsInPlayer {
				// 移除卡牌的逻辑
				gameData.First, done = removeCardFromHand(gameData.First, card)
			}
		}
	case 2:
		for _, cardsInPlayer := range gameData.Sec {
			if card == cardsInPlayer {
				gameData.Sec, done = removeCardFromHand(gameData.Sec, card)
			}
		}
	case 3:
		for _, cardsInPlayer := range gameData.Third {
			if card == cardsInPlayer {
				gameData.Third, done = removeCardFromHand(gameData.Third, card)
			}
		}
	case 4:
		for _, cardsInPlayer := range gameData.Fourth {
			if card == cardsInPlayer {
				gameData.Fourth, done = removeCardFromHand(gameData.Fourth, card)
			}
		}
	}
	if !done {
		slog.Error("[gamelogic/main.go] An error occurred", "error", errors.New("Card Not Found in player"))
		return errors.New("Card Not Found in player")
	}
	if actionData.Turn != player {
		slog.Error("[gamelogic/main.go] An error occurred", "error", errors.New("Not your turn"))
		return errors.New("Not your turn")
	}
	// 更新 LastCard
	gameActionData.LastCard = card
	// 写入数据
	err = WriteGameData(roomID, gameData)
	if err != nil {
		slog.Error("[gamelogic/main.go] An error occurred", "error", err)
		return err
	}
	err = WriteGameActionData(roomID, gameActionData)
	if err != nil {
		slog.Error("[gamelogic/main.go] An error occurred", "error", err)
		return err
	}
	return nil
}
