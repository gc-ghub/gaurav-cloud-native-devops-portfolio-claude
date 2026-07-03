package credly

// badgesResponse mirrors the shape of Credly's unofficial
// /users/{id}/badges.json endpoint (verified live, not a documented API —
// see backend/README.md for the risk this carries).
type badgesResponse struct {
	Data []Badge `json:"data"`
}

type Badge struct {
	IssuedAtDate string `json:"issued_at_date"`
	Issuer       struct {
		Summary  string `json:"summary"`
		Entities []struct {
			Entity struct {
				Name string `json:"name"`
			} `json:"entity"`
		} `json:"entities"`
	} `json:"issuer"`
	BadgeTemplate struct {
		Name     string `json:"name"`
		Url      string `json:"url"`
		ImageURL string `json:"image_url"`
	} `json:"badge_template"`
}

// IssuerName prefers the named entity, falling back to the "issued by ..." summary.
func (b Badge) IssuerName() string {
	if len(b.Issuer.Entities) > 0 && b.Issuer.Entities[0].Entity.Name != "" {
		return b.Issuer.Entities[0].Entity.Name
	}
	return b.Issuer.Summary
}
