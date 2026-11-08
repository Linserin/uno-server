# uno-server
A UNO Game server written in Go

这是 Linserin 练习 Go 的项目，目前正在完善...

路由：
- `/`: 前端开发控制台（静态资源通过 `go:embed` 打包）
- `/ping`: 测试服务器
- `/newRoom`: 创建新房间并发牌（每个玩家发牌 7 张）
- `/fetchPlayerCards`: 获取指定玩家的牌，使用查询方法
- `/fetchAction`: 获取当前对局状态
- `/action`: 出牌；出 ID 为 `1000` 的牌表示跳过（Pass）。使用 Wild / +4 时可附带 `color=1-4` 指定新颜色


## Todo:

已完成：
- 验证玩家已经拥有牌
- `/fetchAction` 使用 poll 来获取对局状态
- 主游戏逻辑 - 数字 & 颜色（Wild / 王牌使用 `color` 字段作为当前生效颜色）
- 主游戏逻辑 - 反转（同色 / 上一张同为反转，两人对战时等同禁止）
- 主游戏逻辑 - +2（同色 & 叠加）
- 主游戏逻辑 - +4（叠加）
- 主游戏逻辑 - 禁止（Skip）
- 主游戏逻辑 - 游戏开始第一张牌不得为 Wild / 王牌
- 跳过（`card=1000`）：不改变上一张牌，只移动回合；有 +2 / +4 罚牌时摸牌并结束回合


等待完成：
- 应该是摸牌 （`/getCardFromRepo`） 之后再交由玩家确认是否出摸出的牌或者直接跳过
- Wild & +4 选择颜色（查询参数中增加 `choseColor`，`int` 类型，值为 1,2,3,4.）
- 游戏胜利判定（Action 数据中增加一个 Winner, 为 1,2,3,4。注意牌库为特殊的玩家 ID; 参见下条 Todo）
- 牌库耗尽后需要从 allCards 复制一个新的切片后删除玩家持有的卡片并洗后作为牌库
- `getCardFromRepo` 写回数据