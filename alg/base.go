package alg

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"log"
	"math/big"
	rand2 "math/rand"
)

var (
	Zero = &big.Int{}
	One  = &big.Int{}
	Two  = &big.Int{}
)

func init() {
	log.Printf("==>> run init() here")
	Zero.SetInt64(0)
	One.SetInt64(1)
	Two.SetInt64(2)
}

// 欧几里得求最大公因数
func GcdEuclid(a, b *big.Int) *big.Int {
	r0, r1 := big.NewInt(0).Set(a), big.NewInt(0).Set(b)
	if a.Cmp(b) == -1 {
		r0.Set(b)
		r1.Set(a)
	}
	r2 := big.NewInt(0)
	r2.Mod(r0, r1)
	for r2.Cmp(Zero) != 0 {
		r0.Set(r1)
		r1.Set(r2)
		r2.Mod(r0, r1)
	}
	return r1
}

// 随机生成nBytes字节序列
func GenRandomBytes(nBytes int) []byte {
	n := (nBytes + 7) / 8
	remain := nBytes % 8
	ret := new(bytes.Buffer)
	for i := 0; i < n; i++ {
		tmp := rand2.Uint64()
		binary.Write(ret, binary.BigEndian, tmp)
	}
	ret2 := ret.Bytes()
	ind := 8 - remain
	if remain == 0 {
		ind = 0
	}
	return ret2[ind:]
}

// 随机生成指定比特个数的大整数，最高位是1
func GenRandomBigInt(nBit int) *big.Int {
	n := (nBit + 7) / 8
	remain := nBit % 8
	randBuf := GenRandomBytes(n)
	if remain > 0 {
		randBuf[0] = (randBuf[0] & ((1 << remain) - 1)) | (1 << (remain - 1))
	} else {
		randBuf[0] = randBuf[0] | 128
	}
	return new(big.Int).SetBytes(randBuf)
}

// 随机生成一个大整数[0,n)=Zn
func GenRandomFromZn(n *big.Int) *big.Int {
	ret, _ := rand.Int(rand.Reader, n)
	return ret
}
func gen_ds(N *big.Int) (*big.Int, int) {
	d, s := new(big.Int).Set(N), 0
	tmp := big.NewInt(0)
	for true {
		tmp.Mod(d, Two)
		if tmp.Cmp(One) == 0 {
			return d, s
		}
		//d=d/2
		d.Set(d.Div(d, Two))
		s = s + 1
	}
	return nil, 0
}

// 判断一个奇数是否为素数
func MillerRabin(N *big.Int, round int) bool {
	if new(big.Int).Mod(N, Two).Cmp(Zero) == 0 {
		//log.Printf("它是偶数")
		return false
	}
	N_1 := new(big.Int).Sub(N, One)
	d, s := gen_ds(N_1)
	for k := 0; k < round; {
		a := GenRandomFromZn(N_1)
		if a.Cmp(Two) == -1 { //a<2
			continue
		}
		x := new(big.Int).Exp(a, d, N)
		y := big.NewInt(0)
		for i := 0; i < s; i++ {
			y.Exp(x, Two, N)
			if y.Cmp(One) == 0 && x.Cmp(One) != 0 && x.Cmp(N_1) != 0 {
				//log.Printf("return false here1 ...")
				return false
			}
			x.Set(y)
		}
		if y.Cmp(One) != 0 {
			//log.Printf("return false here2 y:%+v...", y)
			return false
		}
		k++
	}
	return true
}

func GenPrime(nBit int) (*big.Int, error) {
	if nBit < 3 {
		panic("nBit < 3")
	}
	c := 2
	times := nBit * nBit / c
	for i := 0; i < times; i++ {
		N := GenRandomBigInt(nBit)
		if MillerRabin(N, 40) {
			return N, nil
		}
	}
	return nil, errors.New("not found")
}

type CRT struct {
	N     *big.Int
	pList []*big.Int
}

func GetNPrime(n, nBits int) ([]*big.Int, error) {
	tmp := map[string]struct{}{}
	for i := 0; i < n; {
		x, err := GenPrime(nBits)
		if err != nil {
			continue
		}
		_, ok := tmp[x.String()]
		if !ok {
			tmp[x.String()] = struct{}{}
		} else {
			//log.Printf("%+v has existed", x.String())
			continue
		}
		i++
	}
	res := make([]*big.Int, 0, n)
	for k, _ := range tmp {
		x, fg := new(big.Int).SetString(k, 10)
		if fg {
			res = append(res, x)
		} else {
			return nil, errors.New("new big int error")
		}
	}
	return res, nil
}
func NewCRT(n, nBits int) *CRT {
	pList, err := GetNPrime(n, nBits)
	if err != nil {
		return nil
	}
	N := big.NewInt(1)
	for i := 0; i < n; i++ {
		N.Mul(N, pList[i])
	}
	return &CRT{
		N:     N,
		pList: pList,
	}
}
func NewCRT2(pList []*big.Int) *CRT {
	n := len(pList)
	N := big.NewInt(1)
	for i := 0; i < n; i++ {
		N.Mul(N, pList[i])
	}
	return &CRT{
		N:     N,
		pList: pList,
	}
}
func (c *CRT) N2Zpi(x *big.Int) []*big.Int {
	n := len(c.pList)
	res := make([]*big.Int, 0, n)
	for i := 0; i < n; i++ {
		a := new(big.Int).Mod(x, c.pList[i])
		res = append(res, a)
	}
	return res
}
func MulMod(x, y, P *big.Int) *big.Int {
	tmp := new(big.Int).Mul(x, y)
	tmp.Mod(tmp, P)
	return tmp
}
func (c *CRT) Zpi2N(ai []*big.Int) *big.Int {
	Mi, MiInv := GetMiAndInvFromPList(c.pList)
	n := len(ai)
	if n != len(c.pList) {
		panic("error parameter")
	}
	sum := big.NewInt(0)
	for i := 0; i < n; i++ {
		tmp := MulMod(Mi[i], MiInv[i], c.N)
		tmp.Mul(tmp, ai[i])
		tmp.Mod(tmp, c.N)
		sum.Add(sum, tmp)
		sum.Mod(sum, c.N)
	}
	return sum
}
func (c *CRT) GetMiAndInvFromPList() ([]*big.Int, []*big.Int) {
	return GetMiAndInvFromPList(c.pList)

}
func GetMiAndInvFromPList(pList []*big.Int) ([]*big.Int, []*big.Int) {
	n := len(pList)
	Mi := make([]*big.Int, 0, n)
	for i := 0; i < n; i++ {
		tmp := new(big.Int).SetInt64(1)
		for j := 0; j < n; j++ {
			if j != i {
				tmp.Mul(tmp, pList[j])
			}
		}
		Mi = append(Mi, tmp)
	}
	MiInv := make([]*big.Int, 0, n)
	for i := 0; i < n; i++ {
		tmp := new(big.Int).ModInverse(Mi[i], pList[i])
		MiInv = append(MiInv, tmp)
	}
	return Mi, MiInv
}
