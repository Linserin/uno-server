package gamelogic

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"math/rand"
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
	Winner    int  `json:"winner"`    // 获胜玩家（1-4），0 表示游戏进行中
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
	allCards := AllCards()
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

// refillRepo 牌库耗尽后重建牌库：
// 复制完整的 allCards，删除所有玩家手中持有的牌和当前桌面上的牌，洗牌后作为新牌库。
func refillRepo(gameData *UNO_Cards, lastCard int) {
	unavailable := make(map[int]bool)
	for _, hand := range [][]int{gameData.First, gameData.Sec, gameData.Third, gameData.Fourth} {
		for _, card := range hand {
			unavailable[card] = true
		}
	}
	unavailable[lastCard] = true

	newRepo := make([]int, 0, len(AllCards()))
	for _, card := range AllCards() {
		if !unavailable[card] {
			newRepo = append(newRepo, card)
		}
	}
	rand.Shuffle(len(newRepo), func(i, j int) {
		newRepo[i], newRepo[j] = newRepo[j], newRepo[i]
	})
	gameData.Repo = newRepo
	slog.Info("Deck refilled", "cards", len(newRepo))
}

// drawOneCard 从牌库摸一张牌，牌库为空时先重建牌库
func drawOneCard(gameData *UNO_Cards, lastCard int) (int, error) {
	if len(gameData.Repo) == 0 {
		refillRepo(gameData, lastCard)
	}
	if len(gameData.Repo) == 0 {
		return 0, errors.New("No cards left in deck")
	}
	card := gameData.Repo[0]
	gameData.Repo = gameData.Repo[1:]
	return card, nil
}

// DrawCards 从牌库给指定玩家摸 count 张牌，返回摸到的牌
func DrawCards(gameData *UNO_Cards, player int, count int, lastCard int) ([]int, error) {
	drawn := make([]int, 0, count)
	for i := 0; i < count; i++ {
		card, err := drawOneCard(gameData, lastCard)
		if err != nil {
			return drawn, err
		}
		drawn = append(drawn, card)
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
	return drawn, nil
}

// drawCards 从牌库给指定玩家摸 count 张牌（内部使用）
func drawCards(gameData *UNO_Cards, player int, count int, lastCard int) error {
	_, err := DrawCards(gameData, player, count, lastCard)
	return err
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

// playerHand 返回指定玩家手中的牌
func playerHand(gameData UNO_Cards, player int) []int {
	switch player {
	case 1:
		return gameData.First
	case 2:
		return gameData.Sec
	case 3:
		return gameData.Third
	case 4:
		return gameData.Fourth
	}
	return nil
}

func Logic_Action_main(roomID int, player int, card int, choseColor int) error {
	slog.Debug("[gamelogic/main.go] Card Action.", "roomID", roomID, "player", player, "card", card, "choseColor", choseColor)
	if card == 0 {
		return errors.New("card cannot be 0")
	}
	actionData, err := ReadGameActionData(roomID)
	if err != nil {
		slog.Error("[gamelogic/main.go] An error occurred", "error", err)
		return err
	}
	if actionData.Winner != 0 {
		return errors.New("Game Over")
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
			err = drawCards(&gameData, player, actionData.PlusTwo+actionData.PlusFour, actionData.LastCard)
			if err != nil {
				slog.Error("[gamelogic/main.go] An error occurred", "error", err)
				return err
			}
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

	// 有罚牌待结算时：+2 只能叠加 +2（不要求同色），+4 只能叠加 +4。
	if actionData.PlusTwo > 0 && !IsDrawTwoCard(card) {
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
		if choseColor >= 1 && choseColor <= 4 {
			actionData.Color = choseColor
		} else if actionData.Color < 1 || actionData.Color > 4 {
			actionData.Color = 1
		}
	} else {
		actionData.Color = JustifyCardColor(card)
	}

	// 胜利判定：出完手中所有牌
	if len(playerHand(gameData, player)) == 0 {
		actionData.Winner = player
		slog.Info("Game won", "roomID", roomID, "player", player)
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
