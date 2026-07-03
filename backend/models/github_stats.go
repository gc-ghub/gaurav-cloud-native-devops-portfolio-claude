package models

type GithubStats struct {
	PublicRepos int `json:"public_repos"`
	TotalStars  int `json:"total_stars"`
	TotalForks  int `json:"total_forks"`
}
