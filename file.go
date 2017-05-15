package threeplay

import "encoding/json"

type File struct {
	ID                   uint   `json:"id"`
	ProjectID            uint   `json:"project_id"`
	BatchID              uint   `json:"batch_id"`
	Duration             uint   `json:"duration"`
	Attribute1           string `json:"attribute1"`
	Attribute2           string `json:"attribute2"`
	Attribute3           string `json:"attribute3"`
	VideoID              string `json:"video_id"`
	Name                 string `json:"name"`
	CallbackURL          string `json:"callback_url"`
	Description          string `json:"description"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
	WordCount            uint   `json:"word_count"`
	ThumbnailURL         string `json:"thumbnail_url"`
	LanguageID           int    `json:"language_id"`
	DefaultServiceTypeID int    `json:"default_service_type_id"`
	Downloaded           bool   `json:"downloaded"`
	State                string `json:"state"`
	TurnaroundLevelID    int    `json:"turnaround_level_id"`
	Deadline             string `json:"deadline"`
	BatchName            string `json:"batch_name"`
	ErrorDescription     string `json:"error_description"`
}

type FilesPage struct {
	Files   []File `json:"files"`
	Summary `json:"summary"`
}

type Summary struct {
	CurrentPage  json.Number `json:"current_page"`
	PerPage      json.Number `json:"per_page"`
	TotalEntries json.Number `json:"total_entries"`
	TotalPages   json.Number `json:"total_pages"`
}
