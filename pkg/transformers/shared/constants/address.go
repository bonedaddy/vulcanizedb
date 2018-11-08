package constants

import "github.com/spf13/viper"

func getContractValue(key string, fallback string) string {
	value := viper.GetString(key)
	if value == "" {
		return fallback
	}
	return value
}

var (
	// temporary addresses from Kovan deployment
	CatContractAddress     = getContractValue("contract.cat", "0x2f34f22a00ee4b7a8f8bbc4eaee1658774c624e0")
	DripContractAddress    = getContractValue("contract.drip", "0x891c04639a5edcae088e546fa125b5d7fb6a2b9d")
	FlapperContractAddress = getContractValue("contract.mcd_flap", "0x8868BAd8e74FcA4505676D1B5B21EcC23328d132") // MCD FLAP Contract
	FlipperContractAddress = getContractValue("contract.eth_flip", "0x32D496Ad866D110060866B7125981C73642cc509") // ETH FLIP Contract
	FlopperContractAddress = getContractValue("contract.mcd_flop", "0x6191C9b0086c2eBF92300cC507009b53996FbFFa") // MCD FLOP Contract
	PepContractAddress     = getContractValue("contract.pep", "0xB1997239Cfc3d15578A3a09730f7f84A90BB4975")
	PipContractAddress     = getContractValue("contract.pip", "0x9FfFE440258B79c5d6604001674A4722FfC0f7Bc")
	PitContractAddress     = getContractValue("contract.pit", "0xe7cf3198787c9a4daac73371a38f29aaeeced87e")
	RepContractAddress     = getContractValue("contract.rep", "0xf88bBDc1E2718F8857F30A180076ec38d53cf296")
	VatContractAddress     = getContractValue("contract.vat", "0xcd726790550afcd77e9a7a47e86a3f9010af126b")
	VowContractAddress     = getContractValue("contract.vow", "0x3728e9777B2a0a611ee0F89e00E01044ce4736d1")
)
