package handler

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"
	"unicode/utf8"

	"kurohelper-api/cache"
	"kurohelper-api/dto"

	kurohelperservice "kurohelperservice"
	"kurohelperservice/provider/erogs"

	"github.com/gofiber/fiber/v3"
)

func GetErogsMusic(c fiber.Ctx) error {
	// 快取鍵會跨請求保留，避免引用 Fiber 重用的請求緩衝區，這邊使用Clone來解決。
	keyword := strings.Clone(strings.TrimSpace(c.Query("q")))
	if keyword == "" || utf8.RuneCountInString(keyword) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "q 必須是 1 到 100 字的關鍵字",
			Data:    nil,
		})
	}

	cache.ErogsMusicListStore.Clean()
	list, err := cache.ErogsMusicListStore.Get(keyword)
	if errors.Is(err, kurohelperservice.ErrCacheLost) {
		list, err = erogs.SearchMusicListByKeyword([]string{keyword, kurohelperservice.ZhTwToJp(keyword)})
		if errors.Is(err, kurohelperservice.ErrSearchNoContent) {
			list, err = []erogs.MusicList{}, nil
		}
		if err == nil {
			cache.ErogsMusicListStore.Set(keyword, list)
		}
	}
	if err != nil {
		return respondErogsMusicError(c, err, "GetErogsMusic")
	}

	out := make([]dto.ErogsMusicListResponse, 0, len(list))
	for _, music := range list {
		games := make([]string, 0, len(music.Games))
		image := ""
		for _, game := range music.Games {
			if name := strings.TrimSpace(game.Name); name != "" {
				games = append(games, name)
			}
			if dmm := strings.TrimSpace(game.DMM); image == "" && dmm != "" {
				image = erogs.MakeDMMImageURL(dmm)
			}
		}
		out = append(out, dto.ErogsMusicListResponse{
			ID:           music.ID,
			Name:         music.Name,
			Category:     music.Category,
			AvgTokuten:   music.AvgTokuten,
			TokutenCount: music.TokutenCount,
			Image:        image,
			Games:        games,
		})
	}

	return c.JSON(dto.TResponse[[]dto.ErogsMusicListResponse]{Message: "ok", Data: out})
}

func GetErogsMusicByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(strings.TrimSpace(c.Params("id")))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "id 必須是大於 0 的整數",
			Data:    nil,
		})
	}

	key := strconv.Itoa(id)
	cache.ErogsMusicStore.Clean()
	music, err := cache.ErogsMusicStore.Get(key)
	if errors.Is(err, kurohelperservice.ErrCacheLost) {
		music, err = erogs.SearchMusicByID(id)
		if err == nil && (music == nil || music.ID == 0) {
			err = kurohelperservice.ErrSearchNoContent
		}
		if err == nil {
			cache.ErogsMusicStore.Set(key, music)
		}
	}
	if err != nil {
		return respondErogsMusicError(c, err, "GetErogsMusicByID")
	}

	games := make([]dto.ErogsMusicGameResponse, 0, len(music.GameCategories))
	image := ""
	for _, game := range music.GameCategories {
		if strings.TrimSpace(game.GameName) == "" {
			continue
		}
		gameImage := ""
		if dmm := strings.TrimSpace(game.GameDMM); dmm != "" {
			gameImage = erogs.MakeDMMImageURL(dmm)
		}
		if image == "" {
			image = gameImage
		}
		games = append(games, dto.ErogsMusicGameResponse{
			Name:     game.GameName,
			Model:    game.GameModel,
			Category: game.Category,
			Image:    gameImage,
		})
	}

	playTime, releaseDate := music.PlayTime, music.ReleaseDate
	if playTime == "00:00:00" {
		playTime = ""
	}
	if releaseDate == "0001-01-01" {
		releaseDate = ""
	}

	return c.JSON(dto.TResponse[dto.ErogsMusicResponse]{
		Message: "ok",
		Data: dto.ErogsMusicResponse{
			ID:           music.ID,
			Name:         music.MusicName,
			PlayTime:     playTime,
			ReleaseDate:  releaseDate,
			AvgTokuten:   music.AvgTokuten,
			TokutenCount: music.TokutenCount,
			Image:        image,
			Singers:      splitErogsMusicNames(music.Singers),
			Lyrics:       splitErogsMusicNames(music.Lyrics),
			Compositions: splitErogsMusicNames(music.Compositions),
			Arrangements: splitErogsMusicNames(music.Arrangments),
			Albums:       splitErogsMusicNames(music.Album),
			Games:        games,
		},
	})
}

func splitErogsMusicNames(value string) []string {
	names := make([]string, 0)
	for _, part := range strings.Split(value, ",") {
		if name := strings.TrimSpace(part); name != "" && name != "無" {
			names = append(names, name)
		}
	}
	return names
}

func respondErogsMusicError(c fiber.Ctx, err error, operation string) error {
	status, message := fiber.StatusInternalServerError, "發生錯誤，請稍後再試"
	switch {
	case errors.Is(err, kurohelperservice.ErrSearchNoContent):
		status, message = fiber.StatusNotFound, "找不到音樂"
	case errors.Is(err, kurohelperservice.ErrRateLimit):
		status, message = fiber.StatusTooManyRequests, "批評空間查詢過於頻繁，請稍後再試"
	default:
		slog.Error(operation, "err", err)
	}
	return c.Status(status).JSON(dto.TResponse[any]{Message: message, Data: nil})
}
