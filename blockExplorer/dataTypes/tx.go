package dataTypes

type Tag struct {
	Key   string
	Value string
}

type ResultTx struct {
	Log  string
	Tags []Tag
}

type RawTx struct {
	Height int64
	Index  int    // TODO: Confirm which int type: int, int32 or int64
	Tx     string // TODO: Confirm if this does not contain more than transaction, if it is then []string
	Result ResultTx
}
