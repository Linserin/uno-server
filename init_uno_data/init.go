package init_uno_data

// 初始化房间（发牌）
import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"uno-server/gamelogic"
)

type UNO_Cards struct {
	First  []int `json:"first"`
	Sec    []int `json:"Sec"`
	Third  []int `json:"Third"`
	Fourth []int `json:"Fourth"`
	Repo   []int `json:"repo"`
}

func InitUNO_Mem(roomID int, pep int) error {
	var delivrCard UNO_Cards
	allCards := gamelogic.AllCards()
	// Fxxk the order
	rand.Shuffle(len(allCards), func(i, j int) {
		allCards[i], allCards[j] = allCards[j], allCards[i] // 交换元素
	})
	slog.Info("[initUNO_Mem] Init Data", "roomID", roomID, "players", pep, "cards", allCards)
	// 每一个玩家获取 7 张牌
	var CardCounts int
	var forWhich int = 1
	for _, Card := range allCards {
		slog.Debug("Deliver Card", "player", forWhich, "Card", Card)
		forWhich++
		// 玩家 1-4
		if forWhich == 5 {
			forWhich = 1
		}
		if CardCounts == 28 {
			// 当 4 个玩家手里都有 7 张牌时，将剩余的全部写入牌库（Repo）
			forWhich = 6
		}
		// 分发
		switch forWhich {
		case 1:
			delivrCard.First = append(delivrCard.First, Card)
		case 2:
			delivrCard.Sec = append(delivrCard.Sec, Card)
		case 3:
			delivrCard.Third = append(delivrCard.Third, Card)
		case 4:
			delivrCard.Fourth = append(delivrCard.Fourth, Card)
		case 6:
			delivrCard.Repo = append(delivrCard.Repo, Card)
		}
		// 牌库视为特殊玩家 6 号
		if forWhich != 6 {
			CardCounts++
		}

	}
	// 人数为 2, 3, 4.
	switch pep {
	case 2:
		delivrCard.Repo = append(delivrCard.Repo, delivrCard.Third...)
		delivrCard.Third = []int{}
		delivrCard.Repo = append(delivrCard.Repo, delivrCard.Fourth...)
		delivrCard.Fourth = []int{}
	case 3:
		delivrCard.Repo = append(delivrCard.Repo, delivrCard.Fourth...)
		delivrCard.Fourth = []int{}
	}

	file, err := os.Create(strconv.Itoa(roomID) + ".json")
	if err != nil {
		slog.Error("[init_uno_data/init.go] An error occurred", "error", err)
		return err
	}

	// 从牌库中抽取第一张牌出来
	var ActionData gamelogic.GameActionData
	// 初始化 ActionData
	ActionData.Turn = 1
	ActionData.ClockWise = true
	ActionData.Ban = false
	ActionData.Players = pep
	ActionData.Winner = 0
	// 从牌库抽第一张牌作为起始牌，起始牌不得为 Wild / 王牌
	for {
		ActionData.LastCard = delivrCard.Repo[0]
		delivrCard.Repo = delivrCard.Repo[1:]
		if !gamelogic.IsWildCard(ActionData.LastCard) {
			break
		}
	}
	ActionData.Color = gamelogic.JustifyCardColor(ActionData.LastCard)
	slog.Debug("LastCard Got", "ID", ActionData.LastCard)
	// 序列化并且写入数据
	jsonContent, err := json.Marshal(delivrCard)
	err = os.WriteFile(strconv.Itoa(roomID)+".json", jsonContent, 0644)
	if err != nil {
		slog.Error("An error occurred", "error", err)
		return err
	}
	// 完成 Action
	file.Close()
	jsonContent, err = json.Marshal(ActionData)
	if err != nil {
		slog.Error("An error occurred", "error", err)
		return err
	}
	err = os.WriteFile(strconv.Itoa(roomID)+".data.json", jsonContent, 0644)
	if err != nil {
		slog.Error("An error occurred", "error", err)
		return err
	}
	return nil
}
