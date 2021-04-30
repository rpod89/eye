//package main
package oracle

import (
	//"fmt"
	//"math"

	gecko "github.com/superoo7/go-gecko/v3"
)

// Ask geckoo for the simple price of a coin
func SimplePrice(coin, currency string) (float64, error) {
	cg := gecko.NewClient(nil)
	price, err := cg.SimpleSinglePrice(coin, currency)
	if err != nil {
		return float64(0), err
	}
	return float64(price.MarketPrice), nil
}

/*
func main() {
	price := SimplePrice("bitcoin","usd")			// ask for the price
	round := fmt.Sprintf("%.02f", price)	// format at 2 decimal
	ie754 := math.Float64bits(price)	// print a ieee754 binary version
	fe754 := math.Float64frombits(ie754)	// convert back to float64 version
	fmt.Println("round:", round, "\nie754:", ie754, "\nfe754:", fe754)
}
*/
