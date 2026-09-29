package dto

import "time"

type BrandErogsResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Disband   bool      `json:"disband"`
	GameCount int       `json:"gameCount"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GameErogsResponse struct {
	ID            int       `json:"id"`
	BrandErogsID  int       `json:"brandErogsId"`
	Name          string    `json:"name"`
	Image         string    `json:"image"`
	ImageFromGame bool      `json:"imageFromGame"` // true：顯示圖來自 games（含 erogs 無圖時）
	Category      string    `json:"category"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type ErogsOfficialGameCreatorResponse struct {
	ShubetuType       int    `json:"shubetuType"`
	CreatorName       string `json:"creatorName"`
	ShubetuDetailType int    `json:"shubetuDetailType"`
	ShubetuDetailName string `json:"shubetuDetailName"`
}

type ErogsOfficialGameResponse struct {
	ID                               int                                `json:"id"`
	BrandID                          int                                `json:"brandId"`
	BrandName                        string                             `json:"brandName"`
	Name                             string                             `json:"name"`
	SellDay                          string                             `json:"sellDay"`
	Model                            string                             `json:"model"`
	DMM                              string                             `json:"dmm"`
	Median                           string                             `json:"median"`
	TokutenCount                     string                             `json:"tokutenCount"`
	TotalPlayTimeMedian              string                             `json:"totalPlayTimeMedian"`
	TimeBeforeUnderstandingFunMedian string                             `json:"timeBeforeUnderstandingFunMedian"`
	Okazu                            string                             `json:"okazu"`
	Erogame                          string                             `json:"erogame"`
	Genre                            string                             `json:"genre"`
	BannerURL                        string                             `json:"bannerUrl"`
	SteamID                          string                             `json:"steamId"`
	VndbID                           string                             `json:"vndbId"`
	Shoukai                          string                             `json:"shoukai"`
	Junni                            int                                `json:"junni"`
	Creators                         []ErogsOfficialGameCreatorResponse `json:"creators"`
}
