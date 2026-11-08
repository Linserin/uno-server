package gamelogic

import (
	"errors"
	"log/slog"
)

// DrawCardFromRepo 让玩家从牌库摸一张牌，追加到手牌并写回数据。
// 摸牌后由玩家决定出这张牌（调用 /action）或者直接跳过（/action card=1000）。
func DrawCardFromRepo(roomID int, player int) (int, error) {
	actionData, err := ReadGameActionData(roomID)
	if err != nil {
		slog.Error("[gamelogic/get_card_from_repo.go] An error occurred", "error", err)
		return 0, err
	}
	if actionData.Winner != 0 {
		return 0, errors.New("Game Over")
	}
	if actionData.Turn != player {
		return 0, errors.New("Not your turn")
	}
	if actionData.PlusTwo > 0 || actionData.PlusFour > 0 {
		return 0, errors.New("Must stack or pass")
	}
	if actionData.HasDrawn {
		return 0, errors.New("Already drew a card this turn")
	}

	gameData, err := ReadGameData(roomID)
	if err != nil {
		slog.Error("[gamelogic/get_card_from_repo.go] An error occurred", "error", err)
		return 0, err
	}
	drawn, err := DrawCards(&gameData, player, 1, actionData.LastCard)
	if err != nil {
		slog.Error("[gamelogic/get_card_from_repo.go] An error occurred", "error", err)
		return 0, err
	}
	actionData.HasDrawn = true
	if err := WriteGameData(roomID, gameData); err != nil {
		slog.Error("[gamelogic/get_card_from_repo.go] An error occurred", "error", err)
		return 0, err
	}
	if err := WriteGameActionData(roomID, actionData); err != nil {
		slog.Error("[gamelogic/get_card_from_repo.go] An error occurred", "error", err)
		return 0, err
	}
	slog.Info("Player drew a card", "roomID", roomID, "player", player, "card", drawn[0])
	return drawn[0], nil
}
