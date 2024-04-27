package util_test

import (
	"fmt"
	"math"
	"math/big"

	"github.com/consensus-network/consensusd/util/difficulty"

	"github.com/consensus-network/consensusd/util"
)

func ExampleAmount() {

	a := util.Amount(0)
	fmt.Println("Zero Sompi:", a)

	a = util.Amount(1e8)
	fmt.Println("100,000,000 Sompi:", a)

	a = util.Amount(1e5)
	fmt.Println("100,000 Sompi:", a)
	// Output:
	// Zero Sompi: 0 CSS
	// 100,000,000 Sompi: 1 CSS
	// 100,000 Sompi: 0.001 CSS
}

func ExampleNewAmount() {
	amountOne, err := util.NewAmount(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountOne) //Output 1

	amountFraction, err := util.NewAmount(0.01234567)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountFraction) //Output 2

	amountZero, err := util.NewAmount(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountZero) //Output 3

	amountNaN, err := util.NewAmount(math.NaN())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountNaN) //Output 4

	// Output: 1 CSS
	// 0.01234567 CSS
	// 0 CSS
	// invalid consensus amount
}

func ExampleAmount_unitConversions() {
	amount := util.Amount(44433322211100)

	fmt.Println("Sompi to kCSS:", amount.Format(util.AmountKiloCSS))
	fmt.Println("Sompi to CSS:", amount)
	fmt.Println("Sompi to MilliCSS:", amount.Format(util.AmountMilliCSS))
	fmt.Println("Sompi to MicroCSS:", amount.Format(util.AmountMicroCSS))
	fmt.Println("Sompi to Sompi:", amount.Format(util.AmountSompi))

	// Output:
	// Sompi to kCSS: 444.333222111 kCSS
	// Sompi to CSS: 444333.222111 CSS
	// Sompi to MilliCSS: 444333222.111 mCSS
	// Sompi to MicroCSS: 444333222111 μCSS
	// Sompi to Sompi: 44433322211100 Sompi
}

// This example demonstrates how to convert the compact "bits" in a block header
// which represent the target difficulty to a big integer and display it using
// the typical hex notation.
func ExampleCompactToBig() {
	bits := uint32(419465580)
	targetDifficulty := difficulty.CompactToBig(bits)

	// Display it in hex.
	fmt.Printf("%064x\n", targetDifficulty.Bytes())

	// Output:
	// 0000000000000000896c00000000000000000000000000000000000000000000
}

// This example demonstrates how to convert a target difficulty into the compact
// "bits" in a block header which represent that target difficulty .
func ExampleBigToCompact() {
	// Convert the target difficulty from block 300000 in the bitcoin
	// main chain to compact form.
	t := "0000000000000000896c00000000000000000000000000000000000000000000"
	targetDifficulty, success := new(big.Int).SetString(t, 16)
	if !success {
		fmt.Println("invalid target difficulty")
		return
	}
	bits := difficulty.BigToCompact(targetDifficulty)

	fmt.Println(bits)

	// Output:
	// 419465580
}
