package models

type Certification struct {
	Name     string `json:"name"`
	Issuer   string `json:"issuer"`
	ImageURL string `json:"image_url"`
	Url      string `json:"url"`
	IssuedAt string `json:"issued_at"`
}
