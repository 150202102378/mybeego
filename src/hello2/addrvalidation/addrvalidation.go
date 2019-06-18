package addrvalidation

import (
	"crypto/sha256"
	"math/big"
	"reflect"
	"strings"
)

type AddrValidationUnit struct {
	base58ElementTable string
}

//单例模式，获取验证器实例
var singleAddrValidation *AddrValidationUnit

func GetAddrValidation() *AddrValidationUnit {
	if singleAddrValidation == nil {
		singleAddrValidation = new(AddrValidationUnit)
		singleAddrValidation.base58ElementTable = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	}
	return singleAddrValidation
}

/*
 *参数：string类型的比特币地址，及比特币地址类型（BitoinPrefix，PtSHPrefix，Testnet）
 *功能：对比特币地址进行校验，正确返回true，否则返回false
 */
func (addrvalid *AddrValidationUnit) Verify(addr, kind string) bool {
	//比特币地址种类验证
	switch kind {
	case "BitcoinPrefix":
		if addr[0] != '1' {
			return false
		}
	case "PtSHPrefix":
		if addr[0] != '3' {
			return false
		}
	case "Testnet":
		if addr[0] != 'm' && addr[0] != 'n' {
			return false
		}
	default:
		return false
	}

	//验证地址长度
	if len(addr) < 26 || len(addr) > 35 {
		return false
	}

	//解码
	decoded := addrvalid.deCodeBase58To25Bytes(addr)
	if decoded == nil {
		return false
	}

	hash1 := addrvalid.sha256(decoded[:21])
	hash2 := addrvalid.sha256(hash1)

	return reflect.DeepEqual(decoded[21:], hash2[:4])
}

//计算哈希值
func (addrvalid *AddrValidationUnit) sha256(addr []byte) []byte {
	h := sha256.New()
	h.Write(addr)
	return h.Sum(nil)
}

//进行字符及长度的检验并解码
func (addrvalid *AddrValidationUnit) deCodeBase58To25Bytes(addr string) []byte {
	bigRadix := big.NewInt(58)
	result := big.NewInt(0)
	j := big.NewInt(1)

	for i := len(addr) - 1; i >= 0; i-- {
		//检验字符是否在Base58字母表内
		tmp := strings.IndexAny(addrvalid.base58ElementTable, string(addr[i]))
		if tmp == -1 {
			return nil
		}

		//转化为十进制
		idx := big.NewInt(int64(tmp))
		tmp1 := big.NewInt(0)
		tmp1.Mul(j, idx)

		result.Add(result, tmp1)
		j.Mul(j, bigRadix)
	}

	tmpval := result.Bytes()

	//高位可能为0，要加上去
	var numZeros int
	for numZeros = 0; numZeros < len(addr); numZeros++ {
		if addr[numZeros] != addrvalid.base58ElementTable[0] {
			break
		}
	}
	flen := numZeros + len(tmpval)
	val := make([]byte, flen)
	copy(val[numZeros:], tmpval)

	return val
}
