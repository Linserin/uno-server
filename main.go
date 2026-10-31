package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"uno-server/init_uno_data"
)

// RoomID cardID
// 108 张牌，需要分配？
// 数字牌每色有19张，1~9每色各2张，0每色各1张。由红黄蓝绿4色组成
/* 功能牌
功能牌（Action Card）也由红黄蓝绿4色组成。它们的名称分别为"禁止"、"反转"、"+2"。功能牌每种8张，每色2张，它们有着特殊的功能。
禁止牌（Skip）：打出禁止后，你的下一位家暂停出牌1次，轮到再下家出牌。
反转牌（Reverse）：①在多人对战中：打出反转后，当前出牌时针顺序将反转（由顺时针变为逆时针，或者由逆时针变为顺时针），轮到下家（原来的上家）出牌；②在2人对战中：你的下家（即你的对手）暂停出牌1次，你可以再出一次牌（功能同禁止牌）。
罚抽2张牌（+2）（Draw 2）：打出+2后，你的下家将被罚摸2张牌，并且暂停出牌1次，轮到再下家出牌。 [2]
万能牌
万能牌（Wild Card）在UNO中有很重要的地位，具有很强的功能。它们的牌面是黑色的，且万能牌类没有颜色。名称分别为“变色”和“+4”，每种各4张。万能牌有许多别称，例如黑牌、王牌等。
万能牌（Wild）：改变当前出牌颜色，然后让下一位家出牌。
王牌（+4）（Wild Draw 4）：改变当前出牌颜色，下家被罚摸4张牌并且暂停出牌1次，轮到再下一位家出牌。
补充：+4只有出牌者手中没有与参照牌同颜色(不包括同数字或功能的牌或万能牌)才能合法出牌。
+4质疑（Challenge）：当有玩家打出+4时，被加牌的玩家可以选择“质疑”（认为打出+4的玩家违规使用+4）或“不质疑”（认为打出+4的玩家没有违规使用+4）。
①若被加牌的玩家选择“不质疑”，则被加牌玩家罚抽4张牌，并暂停出牌1次。
②若被加牌的玩家选择“质疑”，打出+4的玩家需要把自己的手牌展示给被加牌玩家看。（在电子端游中，该步骤由系统进行判定，玩家无法查看打出+4的手牌）
a.若质疑成功，即打出+4的玩家违规使用+4，则打出+4的玩家罚抽4张牌，被加牌玩家可以继续出牌。
b.若质疑失败，即打出+4的玩家没有违规使用+4，则被加牌玩家罚抽6张牌，并暂停出牌1次。 [2]
如果玩家们足够守规，且不希望公开自己的牌，也可以不展示手中的牌。 */

/*
数字卡 ID 由三位数组成

红黄蓝绿 对应 100 200 300 400

1-9 对应 N0N 和 N1N（1-9 各色两张); 0 对应 N00

禁止（Skip）：N21 N31

反转（Reverse）：N22 N32

罚两张（Draw 2）： N23 N33

王牌（+4 & Wild）：511-514

Wild（Change Color）：521-524
*/

type API_UnoStru struct {
	RoomID int `json:"roomID"`
	Cardid int `json:"cardID"`
	Role   int `json:"role"`
}

// 从 API 传入的载体变量
var API_UNO_POST_Data API_UnoStru

func main() {
	// 日志初始化
	f, err := os.OpenFile("app.json.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	mw := io.MultiWriter(os.Stdout, f) // 一份数据写两处

	logger := slog.New(slog.NewJSONHandler(mw, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	slog.Info("Hello, UNO Server!!!")
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	http.HandleFunc("/newRoom", func(w http.ResponseWriter, r *http.Request) {
		API_UNO_POST_Data.RoomID = rand.Intn(1000)
		init_uno_data.InitUNO_Mem(API_UNO_POST_Data.RoomID)
		slog.Info("New room generated", "RoomID", API_UNO_POST_Data.RoomID)
		fmt.Fprintln(w, API_UNO_POST_Data.RoomID)
	})

	// 监听
	slog.Info("Ready to listen 8080 port.")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}
