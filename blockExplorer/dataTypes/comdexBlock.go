package dataTypes

type PartSetHeaderComdex struct {
	Total string
	Hash  string
}

type BlockIDComdex struct {
	Hash  string
	Parts PartSetHeaderComdex
}

type HeaderComdex struct {
	Chain_id           string
	Height             string
	Time               string
	Num_txs            string
	Total_txs          string
	Last_block_id      BlockIDComdex
	Last_commit_hash   string
	Data_hash          string
	Validators_hash    string
	NextValidatorsHash string
	Consensus_hash     string
	App_hash           string
	Last_results_hash  string
	Evidence_hash      string
	Proposer_address   string
}

type DataComdex struct{ Txs []string }

type VoteComdex struct {
	Type              int
	Height            string
	Round             string
	Block_id          BlockIDComdex
	Timestamp         string
	Validator_address string
	Validator_index   string
	Signature         string
}

type DuplicateVoteEvidenceComdex struct {
	Pub_key string
	Vote_a  VoteComdex
	Vote_b  VoteComdex
}

type EvidenceDataComdex struct{ Evidence []DuplicateVoteEvidenceComdex }

type CommitComdex struct {
	Block_id   BlockIDComdex
	Precommits []VoteComdex
}

type BlockComdex struct {
	Header      HeaderComdex
	Data        DataComdex
	Evidence    EvidenceDataComdex
	Last_commit CommitComdex
}

type ComdexTx struct {
	Hash   string
	Height int64
}

type ResultComdex struct {
	Block BlockComdex
}

type ComdexBlock struct {
	Result ResultComdex
}
