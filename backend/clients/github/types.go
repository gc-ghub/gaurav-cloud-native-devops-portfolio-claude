package github

// Repo is the subset of GitHub's repo JSON we care about.
type Repo struct {
	Name            string   `json:"name"`
	FullName        string   `json:"full_name"`
	HTMLURL         string   `json:"html_url"`
	Description     string   `json:"description"`
	StargazersCount int      `json:"stargazers_count"`
	ForksCount      int      `json:"forks_count"`
	Language        string   `json:"language"`
	Topics          []string `json:"topics"`
	Fork            bool     `json:"fork"`
	Archived        bool     `json:"archived"`
	UpdatedAt       string   `json:"updated_at"`
}

// graphQLPinnedResponse mirrors the shape of the GitHub GraphQL API's
// pinnedItems query — the only way to read pinned repos, since the REST API
// has no equivalent endpoint.
type graphQLPinnedResponse struct {
	Data struct {
		User struct {
			PinnedItems struct {
				Nodes []struct {
					Name            string `json:"name"`
					Description     string `json:"description"`
					URL             string `json:"url"`
					StargazerCount  int    `json:"stargazerCount"`
					ForkCount       int    `json:"forkCount"`
					PrimaryLanguage struct {
						Name string `json:"name"`
					} `json:"primaryLanguage"`
				} `json:"nodes"`
			} `json:"pinnedItems"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

