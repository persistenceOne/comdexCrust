package dataTypes

// https://tendermint.com/docs/spec/blockchain/blockchain.html
// where there is []byte, string needs to be given to marshall

// type PartSetHeader struct {
// 	Total int32
// 	Hash  string
// }

// type BlockID struct {
// 	Hash  string
// 	Parts PartSetHeader
// }

// type Version struct {
// 	Block uint64
// 	App   uint64
// }

type Header struct {
	// Version            Version
	// Chain_id           string
	Height int64
	// Time               string
	Num_txs   int64
	Total_txs int64
	// Last_block_id      BlockID
	// Last_commit_hash   string
	// Data_hash          string
	// Validators_hash    string
	// NextValidatorsHash string
	// Consensus_hash     string
	// App_hash           string
	// Last_results_hash  string
	// Evidence_hash      string
	// Proposer_address   string
}

// type Data struct{ Txs []string }

// type Vote struct {
// 	Type              byte
// 	Height            int64
// 	Round             int
// 	Block_id          BlockID
// 	Timestamp         string
// 	Validator_address string
// 	Validator_index   int
// 	Signature         string
// }

// type DuplicateVoteEvidence struct {
// 	Pub_key string
// 	Vote_a  Vote
// 	Vote_b  Vote
// }

// type EvidenceData struct{ Evidence []DuplicateVoteEvidence }

// type Commit struct {
// 	Block_id   BlockID
// 	Precommits []Vote
// }

type Block struct {
	Header Header
	// Data        Data
	// Evidence    EvidenceData
	// Last_commit Commit
}

type RawBlock struct {
	Block Block
}
