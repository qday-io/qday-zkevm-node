package etherman

import (
	"time"
	"math/big"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

// Block struct
type Block struct {
	BlockNumber           uint64
	BlockHash             common.Hash
	ParentHash            common.Hash
	GlobalExitRoots       []GlobalExitRoot
	ForcedBatches         []ForcedBatch
	SequencedBatches      [][]SequencedBatch
	VerifiedBatches       []VerifiedBatch
	SequencedForceBatches [][]SequencedForceBatch
	ForkIDs               []ForkID
	ReceivedAt            time.Time
}

// GlobalExitRoot struct
type GlobalExitRoot struct {
	BlockNumber     uint64
	MainnetExitRoot common.Hash
	RollupExitRoot  common.Hash
	GlobalExitRoot  common.Hash
}

// SequencedBatch represents virtual batch
type SequencedBatch struct {
	BatchNumber   uint64
	SequencerAddr common.Address
	TxHash        common.Hash
	Nonce         uint64
	Coinbase      common.Address
	polygonzkevm.PolygonZkEVMBatchData
}

// ForcedBatch represents a ForcedBatch
type ForcedBatch struct {
	BlockNumber       uint64
	ForcedBatchNumber uint64
	Sequencer         common.Address
	GlobalExitRoot    common.Hash
	RawTxsData        []byte
	ForcedAt          time.Time
}

// VerifiedBatch represents a VerifiedBatch
type VerifiedBatch struct {
	BlockNumber uint64
	BatchNumber uint64
	Aggregator  common.Address
	StateRoot   common.Hash
	TxHash      common.Hash
}

// SequencedForceBatch is a sturct to track the ForceSequencedBatches event.
type SequencedForceBatch struct {
	BatchNumber uint64
	Coinbase    common.Address
	TxHash      common.Hash
	Timestamp   time.Time
	Nonce       uint64
	polygonzkevm.PolygonZkEVMForcedBatchData
}

// ForkID is a sturct to track the ForkID event.
type ForkID struct {
	BatchNumber uint64
	ForkID      uint64
	Version     string
}

type EtherHeader struct {
	*types.Header
	BlockHash       common.Hash
}

func (header *EtherHeader) Hash() common.Hash { return header.BlockHash }

type EtherBlock struct {
	*types.Block
	BlockHash       common.Hash
}

func (block *EtherBlock) Hash() common.Hash { return block.BlockHash }

type EthermintBlock struct {
	BlockNumber     *BigInt     `json:"number"`
	BlockHash       common.Hash `json:"hash"`
	BlockParentHash common.Hash `json:"parentHash"`
}

type BigInt struct {
	big.Int
}

func (i *BigInt) UnmarshalJSON(b []byte) error {
	var val string
	err := json.Unmarshal(b, &val)
	if err != nil {
		return err
	}
	_, success := i.SetString(strings.TrimPrefix(val, "0x"), 16)
	if !success {
		return fmt.Errorf("cannot convert string to big integer")
	}
	return nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	// It's negative.
	if number.IsInt64() {
		return rpc.BlockNumber(number.Int64()).String()
	}
	// It's negative and large, which is invalid.
	return fmt.Sprintf("<invalid %d>", number)
}