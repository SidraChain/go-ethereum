package contracts // contracts

import "math/big"

const (
	OwnerContractAddrHex = "0x0000000000000000000000000000000000000010"
	WacContractAddrHex   = "0x0000000000000000000000000000000000000020"
	PacContractAddrHex   = "0x0000000000000000000000000000000000000030"
	SidraTokenAddrHex    = "0x0000000000000000000000000000000000000040"
	MainFaucetAddrHex    = "0x0000000000000000000000000000000000000050"
	WaqfAddrHex          = "0x0000000000000000000000000000000000000060"
	ZakatAddrHex         = "0x0000000000000000000000000000000000000070"
)

var (
	// uint256 largest value (2^256 - 1)
	Uint256Max = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

	Big4 = big.NewInt(4)
	Big5 = big.NewInt(5)
	Big6 = big.NewInt(6)
	Big7 = big.NewInt(7)
)
