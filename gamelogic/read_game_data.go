package gamelogic

import (
	"encoding/json"
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

	var jsonContent UNO_Cards
	err = json.Unmarshal(data, &jsonContent)
	if err != nil {
		slog.Error("[gamelogic/read_game_data.go] An error occurred", "error", err)
		return emptyCard, err
	}
	return jsonContent, nil
}

// 读取 Action 数据
func ReadGameActionData(roomID int) (GameActionData, error) {
	actionData, err := os.ReadFile(strconv.Itoa(roomID) + ".data.json")

	var structedAction GameActionData
	if err != nil {
		return structedAction, err
	}
	err = json.Unmarshal(actionData, &structedAction)
	return structedAction, err
}
