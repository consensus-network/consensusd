package bip32

import "github.com/pkg/errors"

// BitcoinMainnetPrivate is the version that is used for
// bitcoin mainnet bip32 private extended keys.
// Ecnodes to xprv in base58.
var BitcoinMainnetPrivate = [4]byte{
	0x04,
	0x88,
	0xad,
	0xe4,
}

// BitcoinMainnetPublic is the version that is used for
// bitcoin mainnet bip32 public extended keys.
// Ecnodes to xpub in base58.
var BitcoinMainnetPublic = [4]byte{
	0x04,
	0x88,
	0xb2,
	0x1e,
}

// ConsensusMainnetPrivate is the version that is used for
// consensus mainnet bip32 private extended keys.
// Ecnodes to xprv in base58.
var ConsensusMainnetPrivate = [4]byte{
	0x03,
	0x8f,
	0x2e,
	0xf4,
}

// ConsensusMainnetPublic is the version that is used for
// consensus mainnet bip32 public extended keys.
// Ecnodes to kpub in base58.
var ConsensusMainnetPublic = [4]byte{
	0x03,
	0x8f,
	0x33,
	0x2e,
}

// ConsensusTestnetPrivate is the version that is used for
// consensus testnet bip32 public extended keys.
// Ecnodes to ktrv in base58.
var ConsensusTestnetPrivate = [4]byte{
	0x03,
	0x90,
	0x9e,
	0x07,
}

// ConsensusTestnetPublic is the version that is used for
// consensus testnet bip32 public extended keys.
// Ecnodes to ktub in base58.
var ConsensusTestnetPublic = [4]byte{
	0x03,
	0x90,
	0xa2,
	0x41,
}

// ConsensusDevnetPrivate is the version that is used for
// consensus devnet bip32 public extended keys.
// Ecnodes to kdrv in base58.
var ConsensusDevnetPrivate = [4]byte{
	0x03,
	0x8b,
	0x3d,
	0x80,
}

// ConsensusDevnetPublic is the version that is used for
// consensus devnet bip32 public extended keys.
// Ecnodes to xdub in base58.
var ConsensusDevnetPublic = [4]byte{
	0x03,
	0x8b,
	0x41,
	0xba,
}

// ConsensusSimnetPrivate is the version that is used for
// consensus simnet bip32 public extended keys.
// Ecnodes to ksrv in base58.
var ConsensusSimnetPrivate = [4]byte{
	0x03,
	0x90,
	0x42,
	0x42,
}

// ConsensusSimnetPublic is the version that is used for
// consensus simnet bip32 public extended keys.
// Ecnodes to xsub in base58.
var ConsensusSimnetPublic = [4]byte{
	0x03,
	0x90,
	0x46,
	0x7d,
}

func toPublicVersion(version [4]byte) ([4]byte, error) {
	switch version {
	case BitcoinMainnetPrivate:
		return BitcoinMainnetPublic, nil
	case ConsensusMainnetPrivate:
		return ConsensusMainnetPublic, nil
	case ConsensusTestnetPrivate:
		return ConsensusTestnetPublic, nil
	case ConsensusDevnetPrivate:
		return ConsensusDevnetPublic, nil
	case ConsensusSimnetPrivate:
		return ConsensusSimnetPublic, nil
	}

	return [4]byte{}, errors.Errorf("unknown version %x", version)
}

func isPrivateVersion(version [4]byte) bool {
	switch version {
	case BitcoinMainnetPrivate:
		return true
	case ConsensusMainnetPrivate:
		return true
	case ConsensusTestnetPrivate:
		return true
	case ConsensusDevnetPrivate:
		return true
	case ConsensusSimnetPrivate:
		return true
	}

	return false
}
