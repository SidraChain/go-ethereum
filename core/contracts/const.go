package contracts // contracts

import (
	"math/big"

	"github.com/holiman/uint256"
)

const (
	OwnerContractAddrHex     = "0x0000000000000000000000000000000000000010"
	WacContractAddrHex       = "0x0000000000000000000000000000000000000020"
	PacContractAddrHex       = "0x0000000000000000000000000000000000000030"
	RewardDistributorAddrHex = "0x0000000000000000000000000000000000000040"
	MainFaucetAddrHex        = "0x0000000000000000000000000000000000000050"
	WaqfAddrHex              = "0x0000000000000000000000000000000000000060"
	ZakatAddrHex             = "0x0000000000000000000000000000000000000070"
)

var (

	// uint256 max value 32 bytes
	Uint256Max = new(uint256.Int).SetBytes([]byte{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	})

	Big4 = big.NewInt(4)
	Big5 = big.NewInt(5)
	Big6 = big.NewInt(6)
	Big7 = big.NewInt(7)
)
