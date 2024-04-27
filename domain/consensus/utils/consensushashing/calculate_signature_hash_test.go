package consensushashing_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/consensus-network/consensusd/domain/consensus/utils/subnetworks"

	"github.com/kaspanet/go-secp256k1"

	"github.com/consensus-network/consensusd/domain/consensus/utils/consensushashing"
	"github.com/consensus-network/consensusd/domain/consensus/utils/txscript"
	"github.com/consensus-network/consensusd/domain/consensus/utils/utxo"
	"github.com/consensus-network/consensusd/domain/dagconfig"
	"github.com/consensus-network/consensusd/util"

	"github.com/consensus-network/consensusd/domain/consensus/model/externalapi"
)

// shortened versions of SigHash types to fit in single line of test case
const (
	all                = consensushashing.SigHashAll
	none               = consensushashing.SigHashNone
	single             = consensushashing.SigHashSingle
	allAnyoneCanPay    = consensushashing.SigHashAll | consensushashing.SigHashAnyOneCanPay
	noneAnyoneCanPay   = consensushashing.SigHashNone | consensushashing.SigHashAnyOneCanPay
	singleAnyoneCanPay = consensushashing.SigHashSingle | consensushashing.SigHashAnyOneCanPay
)

func modifyOutput(outputIndex int) func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
	return func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
		clone := tx.Clone()
		clone.Outputs[outputIndex].Value = 100
		return clone
	}
}

func modifyInput(inputIndex int) func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
	return func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
		clone := tx.Clone()
		clone.Inputs[inputIndex].PreviousOutpoint.Index = 2
		return clone
	}
}

func modifyAmountSpent(inputIndex int) func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
	return func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
		clone := tx.Clone()
		utxoEntry := clone.Inputs[inputIndex].UTXOEntry
		clone.Inputs[inputIndex].UTXOEntry = utxo.NewUTXOEntry(666, utxoEntry.ScriptPublicKey(), false, 100)
		return clone
	}
}

func modifyScriptPublicKey(inputIndex int) func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
	return func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
		clone := tx.Clone()
		utxoEntry := clone.Inputs[inputIndex].UTXOEntry
		scriptPublicKey := utxoEntry.ScriptPublicKey()
		scriptPublicKey.Script = append(scriptPublicKey.Script, 1, 2, 3)
		clone.Inputs[inputIndex].UTXOEntry = utxo.NewUTXOEntry(utxoEntry.Amount(), scriptPublicKey, false, 100)
		return clone
	}
}

func modifySequence(inputIndex int) func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
	return func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
		clone := tx.Clone()
		clone.Inputs[inputIndex].Sequence = 12345
		return clone
	}
}

func modifyPayload(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
	clone := tx.Clone()
	clone.Payload = []byte{6, 6, 6, 4, 2, 0, 1, 3, 3, 7}
	return clone
}

func modifyGas(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
	clone := tx.Clone()
	clone.Gas = 1234
	return clone
}

func modifySubnetworkID(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
	clone := tx.Clone()
	clone.SubnetworkID = externalapi.DomainSubnetworkID{6, 6, 6, 4, 2, 0, 1, 3, 3, 7}
	return clone
}

func TestCalculateSignatureHashSchnorr(t *testing.T) {
	nativeTx, subnetworkTx, err := generateTxs()
	if err != nil {
		t.Fatalf("Error from generateTxs: %+v", err)
	}

	// Note: Expected values were generated by the same code that they test,
	// As long as those were not verified using 3rd-party code they only check for regression, not correctness
	tests := []struct {
		name                  string
		tx                    *externalapi.DomainTransaction
		hashType              consensushashing.SigHashType
		inputIndex            int
		modificationFunction  func(*externalapi.DomainTransaction) *externalapi.DomainTransaction
		expectedSignatureHash string
	}{
		// native transactions

		// sigHashAll
		{name: "native-all-0", tx: nativeTx, hashType: all, inputIndex: 0,
			expectedSignatureHash: "19d805fc06743bc671a259c6f7cb675a608ef95b3746be0947c3c9b7d8017bb9"},
		{name: "native-all-0-modify-input-1", tx: nativeTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifyInput(1), // should change the hash
			expectedSignatureHash: "5406634b7198d7a70533c38873a0ef7f6bfbc90ee2828bf962ee6cbcee4e4f61"},
		{name: "native-all-0-modify-output-1", tx: nativeTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifyOutput(1), // should change the hash
			expectedSignatureHash: "badf6ddfa426ead1554b37d4fc0079d5f98fdc2e7de79bf8ca8125a11ec327b1"},
		{name: "native-all-0-modify-sequence-1", tx: nativeTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifySequence(1), // should change the hash
			expectedSignatureHash: "9d9dfb63b1bb5a69ff9ff5afee839082cbfc5064951710b35f76aa8cc92eb287"},
		{name: "native-all-anyonecanpay-0", tx: nativeTx, hashType: allAnyoneCanPay, inputIndex: 0,
			expectedSignatureHash: "0b1108b3071c2967fd16f2edf59fcbae5b2bdfedfe5ce0de504cde4f09d1a5f9"},
		{name: "native-all-anyonecanpay-0-modify-input-0", tx: nativeTx, hashType: allAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifyInput(0), // should change the hash
			expectedSignatureHash: "ab99c9f4e208aacdd57150ffb34b3b834fd89ae0d5e8844d7e24a3f36cf8cbbc"},
		{name: "native-all-anyonecanpay-0-modify-input-1", tx: nativeTx, hashType: allAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifyInput(1), // shouldn't change the hash
			expectedSignatureHash: "0b1108b3071c2967fd16f2edf59fcbae5b2bdfedfe5ce0de504cde4f09d1a5f9"},
		{name: "native-all-anyonecanpay-0-modify-sequence", tx: nativeTx, hashType: allAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifySequence(1), // shouldn't change the hash
			expectedSignatureHash: "0b1108b3071c2967fd16f2edf59fcbae5b2bdfedfe5ce0de504cde4f09d1a5f9"},

		// sigHashNone
		{name: "native-none-0", tx: nativeTx, hashType: none, inputIndex: 0,
			expectedSignatureHash: "d117afb7d7a36f522a647ef1d17f35aeac9f2235442b023286f8c89a4d7b7a5b"},
		{name: "native-none-0-modify-output-1", tx: nativeTx, hashType: none, inputIndex: 0,
			modificationFunction:  modifyOutput(1), // shouldn't change the hash
			expectedSignatureHash: "d117afb7d7a36f522a647ef1d17f35aeac9f2235442b023286f8c89a4d7b7a5b"},
		{name: "native-none-0-modify-sequence-0", tx: nativeTx, hashType: none, inputIndex: 0,
			modificationFunction:  modifySequence(0), // should change the hash
			expectedSignatureHash: "1d43d41cb942cb6bd44eff00c060b87e22817f90ef7ebb26b1eaed17bb4a4426"},
		{name: "native-none-0-modify-sequence-1", tx: nativeTx, hashType: none, inputIndex: 0,
			modificationFunction:  modifySequence(1), // shouldn't change the hash
			expectedSignatureHash: "d117afb7d7a36f522a647ef1d17f35aeac9f2235442b023286f8c89a4d7b7a5b"},
		{name: "native-none-anyonecanpay-0", tx: nativeTx, hashType: noneAnyoneCanPay, inputIndex: 0,
			expectedSignatureHash: "a99f566d2712f924e0e783a5054c6720635344637dc965d2ae45879eeb035807"},
		{name: "native-none-anyonecanpay-0-modify-amount-spent", tx: nativeTx, hashType: noneAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifyAmountSpent(0), // should change the hash
			expectedSignatureHash: "677216a420d1a68a59fd40555383648fb7d2a83b5eba3a6c45d21bc8e696f548"},
		{name: "native-none-anyonecanpay-0-modify-script-public-key", tx: nativeTx, hashType: noneAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifyScriptPublicKey(0), // should change the hash
			expectedSignatureHash: "e5c80878d9984265fb334ea6229b5bed7dd09a388bd4c7414538b828c3a2ee24"},

		// sigHashSingle
		{name: "native-single-0", tx: nativeTx, hashType: single, inputIndex: 0,
			expectedSignatureHash: "cc5f2df7d9a62871cda35839abce07e4ce0821c7bf75c9200784250f136d5545"},
		{name: "native-single-0-modify-output-0", tx: nativeTx, hashType: single, inputIndex: 0,
			modificationFunction:  modifyOutput(0), // should change the hash
			expectedSignatureHash: "4b13fa513bf6e6bc95bbede061714788b201c692d15f4a62fbc4045b2bae9199"},
		{name: "native-single-0-modify-output-1", tx: nativeTx, hashType: single, inputIndex: 0,
			modificationFunction:  modifyOutput(1), // shouldn't change the hash
			expectedSignatureHash: "cc5f2df7d9a62871cda35839abce07e4ce0821c7bf75c9200784250f136d5545"},
		{name: "native-single-0-modify-sequence-0", tx: nativeTx, hashType: single, inputIndex: 0,
			modificationFunction:  modifySequence(0), // should change the hash
			expectedSignatureHash: "30858dd27e7da19a4a31cac003151b2f67f02927f667b27b52f04fd445dfd007"},
		{name: "native-single-0-modify-sequence-1", tx: nativeTx, hashType: single, inputIndex: 0,
			modificationFunction:  modifySequence(1), // shouldn't change the hash
			expectedSignatureHash: "cc5f2df7d9a62871cda35839abce07e4ce0821c7bf75c9200784250f136d5545"},
		{name: "native-single-2-no-corresponding-output", tx: nativeTx, hashType: single, inputIndex: 2,
			expectedSignatureHash: "e85a49d337ab82cb4a9aceb2609c1cce27e7d41ee612164c176f9affa97d4190"},
		{name: "native-single-2-no-corresponding-output-modify-output-1", tx: nativeTx, hashType: single, inputIndex: 2,
			modificationFunction:  modifyOutput(1), // shouldn't change the hash
			expectedSignatureHash: "e85a49d337ab82cb4a9aceb2609c1cce27e7d41ee612164c176f9affa97d4190"},
		{name: "native-single-anyonecanpay-0", tx: nativeTx, hashType: singleAnyoneCanPay, inputIndex: 0,
			expectedSignatureHash: "454a20cafaf5f6e22676705d8c3f45eaf8b4a7d70271d381c62b2e8571420847"},
		{name: "native-single-anyonecanpay-2-no-corresponding-output", tx: nativeTx, hashType: singleAnyoneCanPay, inputIndex: 2,
			expectedSignatureHash: "9d34abcdff0bf2eeccba7c8338af2978e3ea8e7dd1f54c3dc7ca897da1393844"},

		// subnetwork transaction
		{name: "subnetwork-all-0", tx: subnetworkTx, hashType: all, inputIndex: 0,
			expectedSignatureHash: "74693821fa095fad4b4960b69cca094693298222bd204d0fc6a31389d7f07b0c"},
		{name: "subnetwork-all-modify-payload", tx: subnetworkTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifyPayload, // should change the hash
			expectedSignatureHash: "440868bc15f5a993607a74ff1bfecf9f4678e16c0d249252d2ff603e9b55c3e3"},
		{name: "subnetwork-all-modify-gas", tx: subnetworkTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifyGas, // should change the hash
			expectedSignatureHash: "41b7a213f9ce2b7ddc773b87e889e2b8c886d2dd5251fc4b27195f812b2eee5c"},
		{name: "subnetwork-all-subnetwork-id", tx: subnetworkTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifySubnetworkID, // should change the hash
			expectedSignatureHash: "25c86b3345c79f63600bc532448ff9faf3b81b923f9a36cca2f783d602a38f3c"},
	}

	for _, test := range tests {
		tx := test.tx
		if test.modificationFunction != nil {
			tx = test.modificationFunction(tx)
		}

		actualSignatureHash, err := consensushashing.CalculateSignatureHashSchnorr(
			tx, test.inputIndex, test.hashType, &consensushashing.SighashReusedValues{})
		if err != nil {
			t.Errorf("%s: Error from CalculateSignatureHashSchnorr: %+v", test.name, err)
			continue
		}

		if actualSignatureHash.String() != test.expectedSignatureHash {
			t.Errorf("%s: expected signature hash: '%s'; but got: '%s'",
				test.name, test.expectedSignatureHash, actualSignatureHash)
		}
	}
}

func TestCalculateSignatureHashECDSA(t *testing.T) {
	nativeTx, subnetworkTx, err := generateTxs()
	if err != nil {
		t.Fatalf("Error from generateTxs: %+v", err)
	}

	// Note: Expected values were generated by the same code that they test,
	// As long as those were not verified using 3rd-party code they only check for regression, not correctness
	tests := []struct {
		name                  string
		tx                    *externalapi.DomainTransaction
		hashType              consensushashing.SigHashType
		inputIndex            int
		modificationFunction  func(*externalapi.DomainTransaction) *externalapi.DomainTransaction
		expectedSignatureHash string
	}{
		// native transactions

		// sigHashAll
		{name: "native-all-0", tx: nativeTx, hashType: all, inputIndex: 0,
			expectedSignatureHash: "2a6cbbc810e0c103f1e0842fab9b333983a7d2c53fd19fdb38a1ea4731576231"},
		{name: "native-all-0-modify-input-1", tx: nativeTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifyInput(1), // should change the hash
			expectedSignatureHash: "f0bcbde70ec77791b7413e164217ae8136e094320070ea5abb6e63823a9ca129"},
		{name: "native-all-0-modify-output-1", tx: nativeTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifyOutput(1), // should change the hash
			expectedSignatureHash: "9fa6cf89cac9a9be7c36f217d5f9210626a518cd600fce693eddfca7ac72bf42"},
		{name: "native-all-0-modify-sequence-1", tx: nativeTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifySequence(1), // should change the hash
			expectedSignatureHash: "1290bd6da54ccfda787a2d3720f7a1a3eea72d338043a1a44aeb6848e9e15586"},
		{name: "native-all-anyonecanpay-0", tx: nativeTx, hashType: allAnyoneCanPay, inputIndex: 0,
			expectedSignatureHash: "36d269f403748f7d297b8f98e4855c46fed0562310d957f6e685f75672f66b53"},
		{name: "native-all-anyonecanpay-0-modify-input-0", tx: nativeTx, hashType: allAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifyInput(0), // should change the hash
			expectedSignatureHash: "cb102ddc585e386e9176536f4c9f786da0226cfb8923c6cbd9e0788d6d0c090f"},
		{name: "native-all-anyonecanpay-0-modify-input-1", tx: nativeTx, hashType: allAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifyInput(1), // shouldn't change the hash
			expectedSignatureHash: "36d269f403748f7d297b8f98e4855c46fed0562310d957f6e685f75672f66b53"},
		{name: "native-all-anyonecanpay-0-modify-sequence", tx: nativeTx, hashType: allAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifySequence(1), // shouldn't change the hash
			expectedSignatureHash: "36d269f403748f7d297b8f98e4855c46fed0562310d957f6e685f75672f66b53"},

		// sigHashNone
		{name: "native-none-0", tx: nativeTx, hashType: none, inputIndex: 0,
			expectedSignatureHash: "39650485cae71858197fc2ae529f784ed2ce9d4f383df97de31984517b55ec33"},
		{name: "native-none-0-modify-output-1", tx: nativeTx, hashType: none, inputIndex: 0,
			modificationFunction:  modifyOutput(1), // shouldn't change the hash
			expectedSignatureHash: "39650485cae71858197fc2ae529f784ed2ce9d4f383df97de31984517b55ec33"},
		{name: "native-none-0-modify-sequence-0", tx: nativeTx, hashType: none, inputIndex: 0,
			modificationFunction:  modifySequence(0), // should change the hash
			expectedSignatureHash: "83058ff5e1c610a3428d064fe836c546195a31c38d5b99182b36b275f21304cd"},
		{name: "native-none-0-modify-sequence-1", tx: nativeTx, hashType: none, inputIndex: 0,
			modificationFunction:  modifySequence(1), // shouldn't change the hash
			expectedSignatureHash: "39650485cae71858197fc2ae529f784ed2ce9d4f383df97de31984517b55ec33"},
		{name: "native-none-anyonecanpay-0", tx: nativeTx, hashType: noneAnyoneCanPay, inputIndex: 0,
			expectedSignatureHash: "da2676f6542a4ee181d54d2ce718a5712739d399a8e81ac6e8d9d3f53a639784"},
		{name: "native-none-anyonecanpay-0-modify-amount-spent", tx: nativeTx, hashType: noneAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifyAmountSpent(0), // should change the hash
			expectedSignatureHash: "f46c8c63f6240cc6f330ed66a66054c13a85bf58428e0b2f48ce498a970225fe"},
		{name: "native-none-anyonecanpay-0-modify-script-public-key", tx: nativeTx, hashType: noneAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifyScriptPublicKey(0), // should change the hash
			expectedSignatureHash: "4e8514d5cb2d02751a5b6ac7ac1a2773d8e2d06ae082b56682a0e3a70f70201a"},

		// sigHashSingle
		{name: "native-single-0", tx: nativeTx, hashType: single, inputIndex: 0,
			expectedSignatureHash: "bf001af851d8fc630495c9441df48d905067c1770222e40ed1b4ca161550598f"},
		{name: "native-single-0-modify-output-0", tx: nativeTx, hashType: single, inputIndex: 0,
			modificationFunction:  modifyOutput(0), // should change the hash
			expectedSignatureHash: "88048698f45e097f4b30bd96ff7d045619277ba2c74eeba2ab9d37c6b3e8618f"},
		{name: "native-single-0-modify-output-1", tx: nativeTx, hashType: single, inputIndex: 0,
			modificationFunction:  modifyOutput(1), // shouldn't change the hash
			expectedSignatureHash: "bf001af851d8fc630495c9441df48d905067c1770222e40ed1b4ca161550598f"},
		{name: "native-single-0-modify-sequence-0", tx: nativeTx, hashType: single, inputIndex: 0,
			modificationFunction:  modifySequence(0), // should change the hash
			expectedSignatureHash: "ab2fb7b5c90a2b02ac9d3899326d0033c731d6653fad64b9a6aae185cdc8042f"},
		{name: "native-single-0-modify-sequence-1", tx: nativeTx, hashType: single, inputIndex: 0,
			modificationFunction:  modifySequence(1), // shouldn't change the hash
			expectedSignatureHash: "bf001af851d8fc630495c9441df48d905067c1770222e40ed1b4ca161550598f"},
		{name: "native-single-2-no-corresponding-output", tx: nativeTx, hashType: single, inputIndex: 2,
			expectedSignatureHash: "a067ffcdb92927116c1c5c68d69375ec3ed3f22fed798b726a9a1b2cf4af32b5"},
		{name: "native-single-2-no-corresponding-output-modify-output-1", tx: nativeTx, hashType: single, inputIndex: 2,
			modificationFunction:  modifyOutput(1), // shouldn't change the hash
			expectedSignatureHash: "a067ffcdb92927116c1c5c68d69375ec3ed3f22fed798b726a9a1b2cf4af32b5"},
		{name: "native-single-anyonecanpay-0", tx: nativeTx, hashType: singleAnyoneCanPay, inputIndex: 0,
			expectedSignatureHash: "0110c75b45a3b7fa78148d7ce8ac0e8d5878b74d5da5715803b67de6eb978b24"},
		{name: "native-single-anyonecanpay-2-no-corresponding-output", tx: nativeTx, hashType: singleAnyoneCanPay, inputIndex: 2,
			expectedSignatureHash: "09154820aa08d7646dd7373500350612e4186421c415590d3b64263df9efca9d"},

		// subnetwork transaction
		{name: "subnetwork-all-0", tx: subnetworkTx, hashType: all, inputIndex: 0,
			expectedSignatureHash: "a98c1f12295a9bc36937c0e53e5ff3aff77eec670031400e18e7a5669095335e"},
		{name: "subnetwork-all-modify-payload", tx: subnetworkTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifyPayload, // should change the hash
			expectedSignatureHash: "6082b1497ae3c0a20ac6996392cc09cacc131ae44a4dd3869b2cfeb702e4aa63"},
		{name: "subnetwork-all-modify-gas", tx: subnetworkTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifyGas, // should change the hash
			expectedSignatureHash: "65047d3f7d9a3cb21c70b227221529a117bb85b4ba5eb6e03ca2f2bf85d6da84"},
		{name: "subnetwork-all-subnetwork-id", tx: subnetworkTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifySubnetworkID, // should change the hash
			expectedSignatureHash: "87deea3ada677132516ca488f559bad282429a53e5683c3fdf7708ec60e0601b"},
	}

	for _, test := range tests {
		tx := test.tx
		if test.modificationFunction != nil {
			tx = test.modificationFunction(tx)
		}

		actualSignatureHash, err := consensushashing.CalculateSignatureHashECDSA(
			tx, test.inputIndex, test.hashType, &consensushashing.SighashReusedValues{})
		if err != nil {
			t.Errorf("%s: Error from CalculateSignatureHashECDSA: %+v", test.name, err)
			continue
		}

		if actualSignatureHash.String() != test.expectedSignatureHash {
			t.Errorf("%s: expected signature hash: '%s'; but got: '%s'",
				test.name, test.expectedSignatureHash, actualSignatureHash)
		}
	}
}

func generateTxs() (nativeTx, subnetworkTx *externalapi.DomainTransaction, err error) {
	genesisCoinbase := dagconfig.SimnetParams.GenesisBlock.Transactions[0]
	genesisCoinbaseTransactionID := consensushashing.TransactionID(genesisCoinbase)

	address1Str := "consensussim:qzpj2cfa9m40w9m2cmr8pvfuqpp32mzzwsuw6ukhfduqpp32mzzws0hk697mh"
	address1, err := util.DecodeAddress(address1Str, util.Bech32PrefixConsensusSim)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding address1: %+v", err)
	}
	address1ToScript, err := txscript.PayToAddrScript(address1)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating script: %+v", err)
	}

	address2Str := "consensussim:qr7w7nqsdnc3zddm6u8s9fex4ysk95hm3v30q353ymuqpp32mzzws0hk697mh"
	address2, err := util.DecodeAddress(address2Str, util.Bech32PrefixConsensusSim)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding address2: %+v", err)
	}
	address2ToScript, err := txscript.PayToAddrScript(address2)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating script: %+v", err)
	}

	txIns := []*externalapi.DomainTransactionInput{
		{
			PreviousOutpoint: *externalapi.NewDomainOutpoint(genesisCoinbaseTransactionID, 0),
			Sequence:         0,
			UTXOEntry:        utxo.NewUTXOEntry(100, address1ToScript, false, 0),
		},
		{
			PreviousOutpoint: *externalapi.NewDomainOutpoint(genesisCoinbaseTransactionID, 1),
			Sequence:         1,
			UTXOEntry:        utxo.NewUTXOEntry(200, address2ToScript, false, 0),
		},
		{
			PreviousOutpoint: *externalapi.NewDomainOutpoint(genesisCoinbaseTransactionID, 2),
			Sequence:         2,
			UTXOEntry:        utxo.NewUTXOEntry(300, address2ToScript, false, 0),
		},
	}

	txOuts := []*externalapi.DomainTransactionOutput{
		{
			Value:           300,
			ScriptPublicKey: address2ToScript,
		},
		{
			Value:           300,
			ScriptPublicKey: address1ToScript,
		},
	}

	nativeTx = &externalapi.DomainTransaction{
		Version:      0,
		Inputs:       txIns,
		Outputs:      txOuts,
		LockTime:     1615462089000,
		SubnetworkID: subnetworks.SubnetworkIDNative,
	}
	subnetworkTx = &externalapi.DomainTransaction{
		Version:      0,
		Inputs:       txIns,
		Outputs:      txOuts,
		LockTime:     1615462089000,
		SubnetworkID: externalapi.DomainSubnetworkID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Gas:          250,
		Payload:      []byte{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
	}

	return nativeTx, subnetworkTx, nil
}

func BenchmarkCalculateSignatureHashSchnorr(b *testing.B) {
	sigHashTypes := []consensushashing.SigHashType{
		consensushashing.SigHashAll,
		consensushashing.SigHashNone,
		consensushashing.SigHashSingle,
		consensushashing.SigHashAll | consensushashing.SigHashAnyOneCanPay,
		consensushashing.SigHashNone | consensushashing.SigHashAnyOneCanPay,
		consensushashing.SigHashSingle | consensushashing.SigHashAnyOneCanPay}

	for _, size := range []int{10, 100, 1000} {
		tx := generateTransaction(b, sigHashTypes, size)

		b.Run(fmt.Sprintf("%d-inputs-and-outputs", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				reusedValues := &consensushashing.SighashReusedValues{}
				for inputIndex := range tx.Inputs {
					sigHashType := sigHashTypes[inputIndex%len(sigHashTypes)]
					_, err := consensushashing.CalculateSignatureHashSchnorr(tx, inputIndex, sigHashType, reusedValues)
					if err != nil {
						b.Fatalf("Error from CalculateSignatureHashSchnorr: %+v", err)
					}
				}
			}
		})
	}
}

func generateTransaction(b *testing.B, sigHashTypes []consensushashing.SigHashType, inputAndOutputSizes int) *externalapi.DomainTransaction {
	sourceScript := getSourceScript(b)
	tx := &externalapi.DomainTransaction{
		Version:      0,
		Inputs:       generateInputs(inputAndOutputSizes, sourceScript),
		Outputs:      generateOutputs(inputAndOutputSizes, sourceScript),
		LockTime:     123456789,
		SubnetworkID: externalapi.DomainSubnetworkID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Gas:          125,
		Payload:      []byte{9, 8, 7, 6, 5, 4, 3, 2, 1},
		Fee:          0,
		Mass:         0,
		ID:           nil,
	}
	signTx(b, tx, sigHashTypes)
	return tx
}

func signTx(b *testing.B, tx *externalapi.DomainTransaction, sigHashTypes []consensushashing.SigHashType) {
	sourceAddressPKStr := "a4d85b7532123e3dd34e58d7ce20895f7ca32349e29b01700bb5a3e72d2570eb"
	privateKeyBytes, err := hex.DecodeString(sourceAddressPKStr)
	if err != nil {
		b.Fatalf("Error parsing private key hex: %+v", err)
	}
	keyPair, err := secp256k1.DeserializeSchnorrPrivateKeyFromSlice(privateKeyBytes)
	if err != nil {
		b.Fatalf("Error deserializing private key: %+v", err)
	}
	for i, txIn := range tx.Inputs {
		signatureScript, err := txscript.SignatureScript(
			tx, i, sigHashTypes[i%len(sigHashTypes)], keyPair, &consensushashing.SighashReusedValues{})
		if err != nil {
			b.Fatalf("Error from SignatureScript: %+v", err)
		}
		txIn.SignatureScript = signatureScript
	}

}

func generateInputs(size int, sourceScript *externalapi.ScriptPublicKey) []*externalapi.DomainTransactionInput {
	inputs := make([]*externalapi.DomainTransactionInput, size)

	for i := 0; i < size; i++ {
		inputs[i] = &externalapi.DomainTransactionInput{
			PreviousOutpoint: *externalapi.NewDomainOutpoint(
				externalapi.NewDomainTransactionIDFromByteArray(&[32]byte{12, 3, 4, 5}), 1),
			SignatureScript: nil,
			Sequence:        uint64(i),
			UTXOEntry:       utxo.NewUTXOEntry(uint64(i), sourceScript, false, 12),
		}
	}

	return inputs
}

func getSourceScript(b *testing.B) *externalapi.ScriptPublicKey {
	sourceAddressStr := "consensussim:qz6f9z6l3x4v3lf9mgf0t934th4nx5kgzu663x9yjh"

	sourceAddress, err := util.DecodeAddress(sourceAddressStr, util.Bech32PrefixConsensusSim)
	if err != nil {
		b.Fatalf("Error from DecodeAddress: %+v", err)
	}

	sourceScript, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		b.Fatalf("Error from PayToAddrScript: %+v", err)
	}
	return sourceScript
}

func generateOutputs(size int, script *externalapi.ScriptPublicKey) []*externalapi.DomainTransactionOutput {
	outputs := make([]*externalapi.DomainTransactionOutput, size)

	for i := 0; i < size; i++ {
		outputs[i] = &externalapi.DomainTransactionOutput{
			Value:           uint64(i),
			ScriptPublicKey: script,
		}
	}

	return outputs
}