package llm

type ReportResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		CoreConclusion string `json:"coreConclusion"`
		Reason         string `json:"reason"`
		Advices        string `json:"advices"`
		TrustIndex     string `json:"trustIndex"`
	} `json:"data"`
}

type AnalysisResult struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}
