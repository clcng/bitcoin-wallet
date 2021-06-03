package keymanager

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/clcng/bitcoin-wallet/pkg/errors"
)

var DefaultRootDerivationPath = DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0}

type DerivationPath []uint32

func (p DerivationPath) toString() string {
	result := "m"
	for _, component := range p {
		var hardened bool
		if component >= 0x80000000 {
			component -= 0x80000000
			hardened = true
		}
		result = fmt.Sprintf("%s/%d", result, component)
		if hardened {
			result += "'"
		}
	}
	return result
}

//  path pattern:
//  m / purpose' / coin' / account' / change / address_index
func ParseDerivationPath(path string) (DerivationPath, error) {
	var result DerivationPath

   components := strings.Split(path, "/")
   if len(components) == 0 {
	   return nil, errors.Coded(1001, "empty derivation path")
   } else if strings.TrimSpace(components[0]) == "" {
	   return nil, errors.Coded(1001, "invalid path")
   } else if strings.TrimSpace(components[0]) == "m" {
	   components = components[1:]
   } else {
		result = append(result, DefaultRootDerivationPath...)
   }

   for _, component := range components {
	component = strings.TrimSpace(component)
	var value uint32

	if strings.HasSuffix(component, "'") {
		value = 0x80000000
		component = strings.TrimSpace(strings.TrimSuffix(component, "'"))
	}

	bigval, ok := new(big.Int).SetString(component, 0)
	if !ok {
		return nil, fmt.Errorf("invalid component: %s", component)
	}
	max := math.MaxUint32 - value
	if bigval.Sign() < 0 || bigval.Cmp(big.NewInt(int64(max))) > 0 {
		if value == 0 {
			return nil, fmt.Errorf("component %v out of allowed range [0, %d]", bigval, max)
		}
		return nil, fmt.Errorf("component %v out of allowed hardened range [0, %d]", bigval, max)
	}
	value += uint32(bigval.Uint64())
	result = append(result, value)
   }
   return result, nil
}
