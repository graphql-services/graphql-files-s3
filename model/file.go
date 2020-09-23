package model

// UploadResponse ...
type UploadResponse struct {
	UID         string `json:"uid"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`

	UploadURL string `json:"uploadURL"`
}

// GetFileResponse
type GetFileResponse struct {
	URL string `json:"url"`
}
