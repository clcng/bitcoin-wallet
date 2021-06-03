package keymanager

import (
	"sync"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

const (
	passphrase = "123456"

	Apostrophe uint32 = 0x80000000 // 0'

	PurposeBIP44 uint32 = 0x8000002C // 44' BIP44
	PurposeBIP49 uint32 = 0x80000031 // 49' BIP49
	PurposeBIP84 uint32 = 0x80000054 // 84' BIP84

	CoinTypeBTC uint32 = 0x80000000
	CoinTypeETH uint32 = 0x8000003c
	CoinTypeCRO uint32 = 0x8000018a
	CoinTypeEOS uint32 = 0x800000c2
)

type KeyManager struct {
	mnemonic   string
	passphrase string
	keys       map[string]*bip32.Key
	mux        sync.Mutex
}

func NewKeyManager(mnemonic string) (*KeyManager, error) {
	km := &KeyManager{
		mnemonic:   mnemonic,
		passphrase: passphrase,
		keys:       make(map[string]*bip32.Key, 0),
	}
	return km, nil
}

func (km *KeyManager) getKey(path string) (*bip32.Key, bool) {
	km.mux.Lock()
	defer km.mux.Unlock()

	key, ok := km.keys[path]
	return key, ok
}

func (km *KeyManager) setKey(path string, key *bip32.Key) {
	km.mux.Lock()
	defer km.mux.Unlock()

	km.keys[path] = key
}

func (km *KeyManager) getMasterKey() (*bip32.Key, error) {
	path := "m"

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	key, err := bip32.NewMasterKey(km.GetSeed())
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)
	return key, nil
}

func (km *KeyManager) GetKey(path string) (*Key, error) {
	key, err := km.GetPathKey(path, 4)
	if err != nil {
		return nil, err
	}

	return &Key{Path: path, Bip32Key: key}, nil
}

//  path pattern:
//  m / purpose' / coin' / account' / change / address_index
func (km *KeyManager) GetPathKey(path string, i int) (*bip32.Key, error) {
	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	if path == "m" {
		return km.getMasterKey()
	}

	derivationPath, err := ParseDerivationPath(path)
	if err != nil {
		return nil, err
	}

	parentPath := derivationPath[:i] 	
	parent, err := km.GetPathKey(parentPath.toString(), i-1)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(derivationPath[i])
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)
	return key, nil
}

func (km *KeyManager) GetSeed() []byte {
	return bip39.NewSeed(km.mnemonic, km.passphrase)
}


