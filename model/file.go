package model

// UploadResponse ...
type UploadResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
	Status      string `json:"status"`

	UploadURL string `json:"uploadURL"`
}

// GetFileResponse
type GetFileResponse struct {
	URL string `json:"url"`
}
