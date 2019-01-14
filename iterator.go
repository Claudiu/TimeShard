package timeshard

type Iterator interface {
	Init()

	HasNext() bool

	Meta() []uint64
	GetMeta(MetaType) uint64
	Value() []byte
}
