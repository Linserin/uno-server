package gamelogic

import (
	"encoding/json"
	"os"
	"strconv"
)

func WriteGameData(roomID int, data UNO_Cards) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(strconv.Itoa(roomID)+".json", jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func WriteGameActionData(roomID int, data GameActionData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(strconv.Itoa(roomID)+".data.json", jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}
