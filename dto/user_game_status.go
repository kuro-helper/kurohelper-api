package dto

const (
	UserGameStatusNone     = 0 // 沒有狀態
	UserGameStatusFinished = 1 // 遊玩完畢
	UserGameStatusPlaying  = 2 // 遊玩中
	UserGameStatusStalled  = 3 // 暫停遊玩
	UserGameStatusDropped  = 4 // 放棄遊玩
)

func ValidUserGameStatus(status int) bool {
	return status >= UserGameStatusNone && status <= UserGameStatusDropped
}
