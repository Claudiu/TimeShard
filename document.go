package timeshard

type Document struct {
	Title      string            `json:"title"`
	Operations Snapshot          `json:"ops"`
	Meta       map[string]string `json:"meta"`
}