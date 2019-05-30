package timeshard

type Document struct {
	Operations Snapshot `json:"ops"`
}
