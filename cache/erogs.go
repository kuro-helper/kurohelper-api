package cache

import (
	"time"

	"kurohelperservice/provider/erogs"
)

var cacheLostTime = 4 * time.Hour

var ErogsGameStore = NewStore[*erogs.Game](cacheLostTime)
