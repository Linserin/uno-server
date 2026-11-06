package gamelogic

import (
	"errors"
	"log/slog"
)

func getCardFromRepo(roomID int, player int) (int, error) {
	var gameData UNO_Cards
	gameData, err := ReadGameData(roomID)
	if err != nil {
		slog.Error("[gamelogic/get_card_from_repo.go] An error occurred", "error", err)
		return 0, err
	}
	var card int
	switch player {
	case 1:
		card = gameData.Repo[0]
		gameData.Repo = gameData.Repo[1:] // Remove the first card from the repository
	case 2:
		card = gameData.Repo[0]
		gameData.Repo = gameData.Repo[1:]
	case 3:
		card = gameData.Repo[0]
		gameData.Repo = gameData.Repo[1:]
	case 4:
		card = gameData.Repo[0]
		gameData.Repo = gameData.Repo[1:]
	default:
		slog.Error("[gamelogic/get_card_from_repo.go] An error occurred", "error", errors.New("Invalid player number"))
		return 0, errors.New("Invalid player number")
	}
	return card, nil
}
