package init_uno_data

// 初始化房间（发牌）
import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type UNO_Cards struct {
	First  []int `json:"first"`
	Sec    []int `json:"Sec"`
	Third  []int `json:"Third"`
	Fourth []int `json:"Fourth"`
	Repo   []int `json:"repo"`
}

func InitUNO_Mem(roomID int) {
	var delivrCard UNO_Cards
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
	// Fxxk the order
	// This might be an old method
	rand.Seed(time.Now().UnixNano()) // 初始化随机种子
	rand.Shuffle(len(allCards), func(i, j int) {
		allCards[i], allCards[j] = allCards[j], allCards[i] // 交换元素
	})
	fmt.Println(allCards)
	slog.Info("[initUNO_Mem] Init Data", "roomID", "roomID")
	// 每一个玩家获取 7 张牌
	var CardCounts int
	var forWhich int = 1
	for _, Card := range allCards {
		slog.Info("Deliver Card", "player", forWhich, "Card", Card)
		forWhich++
		// 玩家 1-4
		if forWhich == 5 {
			forWhich = 1
		}
		if CardCounts == 49 {
			// 当每个玩家手里都有 7 张牌时，将剩余的全部写入牌库（Repo）
			forWhich = 6
		}
		// 牌库视为特殊玩家 6 号
		if forWhich != 6 {
			CardCounts++
		}

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

	}

	file, err := os.Create(strconv.Itoa(roomID) + ".txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonContent, err := json.Marshal(delivrCard)
	err = os.WriteFile(strconv.Itoa(roomID)+".txt", jsonContent, 0644)
	if err != nil {
		slog.Error("An error occurred", "error", err)
	}
	defer file.Close()
}
