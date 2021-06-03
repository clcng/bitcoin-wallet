package server

type GenerateHDAddressInput struct {
	Mnemonic string `json:"mnemonic"`
	Path string `json:"path"`
}