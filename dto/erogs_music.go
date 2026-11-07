package dto

type ErogsMusicListResponse struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	AvgTokuten   float64  `json:"avgTokuten"`
	TokutenCount int      `json:"tokutenCount"`
	Image        string   `json:"image"`
	Games        []string `json:"games"`
}

type ErogsMusicGameResponse struct {
	Name     string `json:"name"`
	Model    string `json:"model"`
	Category string `json:"category"`
	Image    string `json:"image"`
}

type ErogsMusicResponse struct {
	ID           int                      `json:"id"`
	Name         string                   `json:"name"`
	PlayTime     string                   `json:"playTime"`
	ReleaseDate  string                   `json:"releaseDate"`
	AvgTokuten   float64                  `json:"avgTokuten"`
	TokutenCount int                      `json:"tokutenCount"`
	Image        string                   `json:"image"`
	Singers      []string                 `json:"singers"`
	Lyrics       []string                 `json:"lyrics"`
	Compositions []string                 `json:"compositions"`
	Arrangements []string                 `json:"arrangements"`
	Albums       []string                 `json:"albums"`
	Games        []ErogsMusicGameResponse `json:"games"`
}
