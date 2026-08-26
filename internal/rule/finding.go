package rule

type Finding struct {
	Severity       string
	Rule           string
	Path           string
	Message        string
	Recommendation string
}
