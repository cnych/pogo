package parser

type VideoInfo struct {
	Site string `json:"site"`
	Title string `json:"title"`
	Duration int64 `json:"duration"`
	Url string `json:"url"`
	DownloadInfo map[string]interface{} `json:"download_info"`
}
