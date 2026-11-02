package gamelogic

import (
	"log/slog"
)

type UNO_Cards struct {
	First  []int `json:"first"`
	Sec    []int `json:"Sec"`
	Third  []int `json:"Third"`
	Fourth []int `json:"Fourth"`
	Repo   []int `json:"repo"`
}

func Player_Get_Their_Card(roomID int, player int) ([]int, error) {
	slog.Debug("[gamelogic/gamelogic_fetch.go] Ready to read a game file", "file", roomID, "player", player)
	gameData, err := ReadGameData(roomID)
	if err != nil {
		slog.Error("[gamelogic/read_game_data.go] An error occurred", "error", err)
		return []int{}, err
	}
	switch player {
	case 1:
		return gameData.First, nil
	case 2:
		return gameData.Sec, nil
	case 3:
		return gameData.Third, nil
	case 4:
		return gameData.Fourth, nil
	}
	return []int{}, nil
}
