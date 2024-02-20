package contracts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// OwnerContractAddr is the address of the owner contract
	OwnerContractAddr = common.HexToAddress(OwnerContractAddrHex)
	// WacContractAddr is the address of the WAC contract
	WacContractAddr = common.HexToAddress(WacContractAddrHex)
	// PacContractAddr is the address of the PAC contract
	PacContractAddr = common.HexToAddress(PacContractAddrHex)
	// SidraTokenAddr is the address of the Sidra token contract
	SidraTokenAddr = common.HexToAddress(SidraTokenAddrHex)
	// MainFaucetAddr is the address of the main faucet contract
	MainFaucetAddr = common.HexToAddress(MainFaucetAddrHex)
	// WaqfAddr is the address of the waqf contract
	WaqfAddr = common.HexToAddress(WaqfAddrHex)
	// ZakatAddr is the address of the zakat contract
	ZakatAddr = common.HexToAddress(ZakatAddrHex)

	SystemWallets = map[common.Address]bool{
		OwnerContractAddr: true,
		WacContractAddr:   true,
		PacContractAddr:   true,
		SidraTokenAddr:    true,
		MainFaucetAddr:    true,
		WaqfAddr:          true,
		ZakatAddr:         true,
	}
)

// https://docs.soliditylang.org/en/latest/internals/layout_in_storage.html#mappings-and-dynamic-arrays

func ComputeMappingHash(addr *common.Address, slot *big.Int) common.Hash {
	// Convert the slot number to 32 bytes hash with leading zeros
	p := common.BytesToHash(slot.Bytes())

	// Convert the 20 bytes address to 32 bytes hash with leading zeros
	hK := common.BytesToHash(addr.Bytes())

	// Concatenate the key and slot number and convert it to bytes array of 64 bytes
	concatenated := append(hK.Bytes(), p.Bytes()...)

	// Compute and return the Keccak-256 hash of the concatenated key and slot number
	return crypto.Keccak256Hash(concatenated)
}

func NumOfPoolAtAddr(addr *common.Address, statedb *state.StateDB) (*big.Int, int) {
	// Get the state of the WalletAccessControl contract.
	pacState := statedb.GetOrNewStateObject(PacContractAddr)
	// Calculate the keccak256 hash of the key and slot number.
	// INFO: The slot number is from storage layout of the PoolAccessControl contract.
	keyHash := ComputeMappingHash(addr, common.Big1)
	// Get the value of the key from the state which is length of the array.
	value := pacState.GetState(keyHash).Big()
	// Calculate the keccak256 hash of the key which is the first element of the array.
	key := crypto.Keccak256Hash(keyHash.Bytes()).Big()
	// Return the key and value.
	return key, int(value.Int64())
}

func GetListOfPoolAtAddr(start *big.Int, value int, statedb *state.StateDB) []big.Int {
	// Get the state of the WalletAccessControl contract.
	pacState := statedb.GetOrNewStateObject(PacContractAddr)

	list := make([]big.Int, 0, value)
	for i := 0; i < value; i++ {
		key := big.NewInt(0).Add(start, big.NewInt(int64(i)))
		// Get the value of the key from the state.
		v := pacState.GetState(common.BigToHash(key)).Big()
		list = append(list, *v)
	}
	return list
}

func IsInSamePool(addr1 *common.Address, addr2 *common.Address, statedb *state.StateDB) bool {
	keyAddr1, valueAddr1 := NumOfPoolAtAddr(addr1, statedb)
	if valueAddr1 == 0 {
		return false
	}
	keyAddr2, valueAddr2 := NumOfPoolAtAddr(addr2, statedb)
	if valueAddr2 == 0 {
		return false
	}
	list1 := GetListOfPoolAtAddr(keyAddr1, valueAddr1, statedb)
	list2 := GetListOfPoolAtAddr(keyAddr2, valueAddr2, statedb)
	for _, v1 := range list1 {
		for _, v2 := range list2 {
			if v1.Cmp(&v2) == 0 {
				return true
			}
		}
	}
	return false
}

func WalletStatus(addr *common.Address, statedb *state.StateDB) *big.Int {
	if IsSystemContract(*addr) {
		// Return 1 if the address is nil or one of the system wallets.
		return common.Big1
	}
	if addr == nil {
		// Return 0 if the address is nil.
		return common.Big0
	}
	// Get the state of the WalletAccessControl contract.
	wacState := statedb.GetOrNewStateObject(WacContractAddr)

	// Calculate the keccak256 hash of the key and slot number.
	// INFO: The slot number is from storage layout of the WalletAccessControl contract.
	keyHash := ComputeMappingHash(addr, Big6)

	// Get the value of the key from the state.
	value := wacState.GetState(keyHash).Big()

	return value
}

func IsSystemContract(addr common.Address) bool {
	return SystemWallets[addr]
}
func InBlackList(value *big.Int) bool {
	return value.Cmp(common.Big2) == 0
}
func InWhiteListOrUnlisted(value *big.Int) bool {
	return value.Cmp(common.Big0) == 0 || value.Cmp(common.Big1) == 0
}
func InGreyList(value *big.Int) bool {
	// greater than 1
	return value.Cmp(common.Big2) == 1
}
func InSendingGreyList(value *big.Int) bool {
	return value.Cmp(common.Big3) == 0
}
func InRecievingGreyList(value *big.Int) bool {
	return value.Cmp(Big4) == 0
}
func InPoolGreyList(value *big.Int) bool {
	return value.Cmp(Big5) == 0
}

func IsTransactionAllowed(tx *types.Transaction, sender *common.Address, statedb *state.StateDB) bool {
	// Get the state of the WalletAccessControl contract.
	recipient := tx.To()
	senderStatus := WalletStatus(sender, statedb)

	// If the sender is whitelisted and the receiver is nil, return true.
	// This is to allow the creation of new contracts.
	if InWhiteListOrUnlisted(senderStatus) && recipient == nil {
		// Return true if the sender is whitelisted and the receiver is nil.
		return true
	}
	// Get the state of the receiver.
	receiverStatus := WalletStatus(recipient, statedb)
	if InWhiteListOrUnlisted(senderStatus) && InWhiteListOrUnlisted(receiverStatus) {
		// Return true if both sender and receiver are whitelisted.
		return true
	}
	if InBlackList(senderStatus) || InBlackList(receiverStatus) {
		// Return false if either sender or receiver is blacklisted.
		return false
	}
	if InGreyList(senderStatus) && IsSystemContract(*recipient) {
		// Return true if the sender is greylisted and the receiver is one of the system wallets.
		return true
	}
	if InPoolGreyList(senderStatus) || InPoolGreyList(receiverStatus) {
		return IsInSamePool(sender, recipient, statedb)
	}
	if !InSendingGreyList(senderStatus) && !InRecievingGreyList(receiverStatus) {
		// Return true if the sender is not greylisted for sending and the receiver is not greylisted for receiving.
		return true
	}
	// Return false if none of the above conditions are met.
	return false
}
