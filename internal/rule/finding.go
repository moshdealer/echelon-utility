package rule

type Finding struct {
	Severity       string `json:"severity"`
	Rule           string `json:"rule"`
	Source         string `json:"source,omitempty"`
	Path           string `json:"path"`
	Message        string `json:"message"`
	Recommendation string `json:"recommendation"`
}
