package gamelogic

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
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
	file, err := os.Open(strconv.Itoa(roomID) + ".json")
	if err != nil {
		slog.Error("[gamelogic/fetch.go] An error occurred", "error", err)
		return []int{}, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		slog.Error("[gamelogic/fetch.go] An error occurred", "error", err)
		return []int{}, err
	}
	fmt.Println(string(data))

	var jsonContent UNO_Cards
	err = json.Unmarshal(data, &jsonContent)
	if err != nil {
		slog.Error("[gamelogic/fetch.go] An error occurred", "error", err)
		return []int{}, err
	}
	fmt.Println(jsonContent)
	switch player {
	case 1:
		return jsonContent.First, nil
	case 2:
		return jsonContent.Sec, nil
	case 3:
		return jsonContent.Third, nil
	case 4:
		return jsonContent.Fourth, nil
	}
	return []int{}, nil
}
