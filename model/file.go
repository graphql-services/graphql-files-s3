package model

// File ...
type File struct {
	UID         string `json:"uid"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
	URL         string `json:"url"`
	Name        string `json:"name"`
}
