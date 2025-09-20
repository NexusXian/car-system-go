package response

type DriverClassificationResult struct {
	ClassificationTime string `json:"classificationTime"`
	TotalDrivers       int    `json:"totalDrivers"`
	Categories         []struct {
		Type  string `json:"type"`
		Count int    `json:"count"`
	} `json:"categories"`
}
