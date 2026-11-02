package gamelogic

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
)

func ReadGameData(roomID int) (UNO_Cards, error) {
	var emptyCard UNO_Cards
	slog.Debug("[gamelogic/read_game_data.go] Ready to read a game file", "file", roomID)
	file, err := os.Open(strconv.Itoa(roomID) + ".json")
	if err != nil {
		slog.Error("[gamelogic/read_game_data.go] An error occurred", "error", err)
		return emptyCard, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		slog.Error("[gamelogic/read_game_data.go] An error occurred", "error", err)
		return emptyCard, err
	}
	fmt.Println(string(data))

	var jsonContent UNO_Cards
	err = json.Unmarshal(data, &jsonContent)
	if err != nil {
		slog.Error("[gamelogic/read_game_data.go] An error occurred", "error", err)
		return emptyCard, err
	}
	return jsonContent, nil
}
