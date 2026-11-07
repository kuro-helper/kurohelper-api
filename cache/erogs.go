package cache

import (
	"time"

	"kurohelperservice/provider/erogs"
)

var cacheLostTime = 4 * time.Hour

var ErogsGameStore = NewStore[*erogs.Game](cacheLostTime)

var ErogsMusicListStore = NewStore[[]erogs.MusicList](5 * time.Minute)
var ErogsMusicStore = NewStore[*erogs.Music](cacheLostTime)
