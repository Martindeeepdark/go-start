package defs

// TxOptions represents transaction options
type TxOptions struct {
	Isolation int
	ReadOnly  bool
}

// Stats represents database statistics
type Stats struct {
	MaxOpenConnections int
	OpenConnections    int
	InUse              int
	Idle               int
}
