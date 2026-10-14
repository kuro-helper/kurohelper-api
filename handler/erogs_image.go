package handler

import "strings"

// resolveErogsDisplayImage erogs 有圖用 erogs；否則用 games 備選圖。
// imageFromGame 表示顯示路徑屬 games（含 erogs 無圖、備選也為空的情況）。
func resolveErogsDisplayImage(erogsImage, gameImageURL string) (image string, imageFromGame bool) {
	if img := strings.TrimSpace(erogsImage); img != "" {
		return img, false
	}
	return strings.TrimSpace(gameImageURL), true
}
