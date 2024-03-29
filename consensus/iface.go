/**
 * Copyright 2022 Grid 7, LLC (DBA Taekion)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * ------------------------------------------------------------------------------
 */

package consensus

// ConsensusService provides methods that allow the consensus engine to issue commands and requests
// back to the validator.
type ConsensusService interface {
	// -- P2P --

	// SendTo sends a consensus message to a specific, connected peer.
	SendTo(peer PeerId, messageType string, payload []byte) error
	// Broadcast broadcasts a message to all connected peers.
	Broadcast(messageType string, payload []byte) error

	// -- Block Creation --

	// InitializeBlock initializes a new block built on the block with the given previous id and
	// begins adding batches to it. If no previous id is specified, the current
	// head will be used.
	InitializeBlock(previousId BlockId) error
	// SummarizeBlock stops adding batches to the current block and returns a summary of its
	// contents.
	SummarizeBlock() ([]byte, error)
	// FinalizeBlock inserts the given consensus data into the block and signs it. If this call is successful, the
	// consensus engine will receive the block afterwards.
	FinalizeBlock(data []byte) (BlockId, error)
	// CancelBlock stops adding batches to the current block and abandons it.
	CancelBlock() error

	// -- Block Directives --

	// CheckBlocks updates the prioritization of blocks to check.
	CheckBlocks(blockIds []BlockId) error
	// CommitBlock updates the block that should be committed.
	CommitBlock(blockId BlockId) error
	// IgnoreBlock signals that this block is no longer being committed.
	IgnoreBlock(blockId BlockId) error
	// FailBlock marks this block as invalid from the perspective of consensus.
	FailBlock(blockId BlockId) error

	// -- Queries --

	// GetBlocks retrieves consensus-related information about blocks.
	GetBlocks(blockIds []BlockId) (map[BlockId]Block, error)
	// GetChainHead gets the chain head block.
	GetChainHead() (Block, error)
	// GetSettings reads the value of settings as of the given block.
	GetSettings(blockId BlockId, keys []string) (map[string]string, error)
	// GetState reads values in state as of the given block.
	GetState(blockId BlockId, addresses []string) (map[string][]byte, error)
}

// ConsensusEngineImpl must be implemented by any ConsensusEngine implementation.
type ConsensusEngineImpl interface {
	// Name gets the name of the engine, typically the algorithm being implemented.
	Name() string

	// Version gets the version of this engine.
	Version() string

	// Start is called after the engine is initialized, when a connection to the validator has been
	// established. Requests to the validator are sent via calls to the interface methods on service.
	// Updates from the validator are received via updateChan.
	Start(startupState StartupState, service ConsensusService, updateChan chan ConsensusUpdate) error
}
