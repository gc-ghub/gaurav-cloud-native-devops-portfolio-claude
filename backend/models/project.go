package models

// Project mirrors the fields ProjectCard.svelte destructures via $props(),
// so the frontend needs no remapping of the backend's JSON response.
type Project struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tech        []string `json:"tech"`
	Href        string   `json:"href"`
	Stars       int      `json:"stars"`
	Forks       int      `json:"forks"`
}
