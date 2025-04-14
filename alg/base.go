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
	Zero     = &big.Int{}
	One      = &big.Int{}
	Two      = &big.Int{}
	Three    = &big.Int{}
	Four     = &big.Int{}
	NegOne   = &big.Int{}
	BigPrime *big.Int
)

func init() {
	log.Printf("==>> run init() here")
	Zero.SetInt64(0)
	One.SetInt64(1)
	Two.SetInt64(2)
	Three.SetInt64(3)
	Four.SetInt64(4)
	NegOne.SetInt64(-1)
	BigPrime, _ = new(big.Int).SetString("186151550718860665121779391551741013053", 10)
}

// 欧几里得求最大公因数
func GcdEuclid(a, b *big.Int) *big.Int {
	r0, r1 := new(big.Int).Set(a), new(big.Int).Set(b)
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
func RandomBytes(nBytes int) []byte {
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
func RandomBigInt(nBit int) *big.Int {
	n := (nBit + 7) / 8
	remain := nBit % 8
	randBuf := RandomBytes(n)
	if remain > 0 {
		randBuf[0] = (randBuf[0] & ((1 << remain) - 1)) | (1 << (remain - 1))
	} else {
		randBuf[0] = randBuf[0] | 128
	}
	return new(big.Int).SetBytes(randBuf)
}

// 随机生成一个大整数[0,n-1]=Zn
func RandomFromZn(n *big.Int) *big.Int {
	ret, _ := rand.Int(rand.Reader, n)
	return ret
}

// 随机生成一个大整数[1,n-1]=Zn
func RandomFromZnNotZero2(n *big.Int) *big.Int {
	n1 := new(big.Int).Sub(n, One)
	ret, _ := rand.Int(rand.Reader, n1)
	if ret.Cmp(Zero) == 0 {
		return n1
	}
	return ret
}
func RandomFromZnNotZero(n *big.Int) *big.Int {
	for true {
		ret, _ := rand.Int(rand.Reader, n)
		if ret.Cmp(Zero) != 0 {
			return ret
		}
	}
	return nil
}

/*
将一个整数分解成以下形式: N=d*2^s,d是奇数，s是2的幂
*/
func NTods(N *big.Int) (*big.Int, int) {
	d, s := new(big.Int).Set(N), 0
	tmp := big.NewInt(0)
	for true {
		tmp.Mod(d, Two)
		if tmp.Cmp(One) == 0 {
			return d, s
		}
		//d=d/2
		d.Set(d.Div(d, Two))
		//d.Div(d, Two)
		s = s + 1
	}
	return nil, 0
}

// 判断一个奇数是否为素数
func MillerRabin(N *big.Int, round int) bool {
	tmp := new(big.Int)
	if tmp.Mod(N, Two).Cmp(Zero) == 0 {
		return false
	}
	NMinusOne := tmp.Sub(N, One)
	d, s := NTods(NMinusOne)
	y := big.NewInt(0)
	for k := 0; k < round; {
		aaa := RandomFromZnNotZero(NMinusOne)
		if aaa.Cmp(Two) == -1 { //a<2
			continue
		}
		x := aaa.Exp(aaa, d, N)
		y.SetInt64(0)
		for i := 0; i < s; i++ {
			y.Exp(x, Two, N)
			if y.Cmp(One) == 0 && x.Cmp(One) != 0 && x.Cmp(NMinusOne) != 0 {
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

// 生成一个大素数，位数是nBit
func GenPrime(nBit int) (*big.Int, error) {
	if nBit < 3 {
		panic("nBit < 3")
	}
	c := 2
	times := nBit * nBit / c
	for i := 0; i < times; i++ {
		N := RandomBigInt(nBit)
		if MillerRabin(N, 20) {
			return N, nil
		}
	}
	return nil, errors.New("not found")
}

// 随机生成 n 个 nBits 长的素数
func GenPrimeList(n, nBits int) ([]*big.Int, error) {
	mList, err := genNPrime(n, nBits)
	if err != nil {
		return nil, errors.New("getNPrime error")
	}
	return mList, nil
}
func MulBigIntList(mList []*big.Int) *big.Int {
	if len(mList) == 0 {
		panic("input is null")
	}
	N := big.NewInt(1)
	for i := 0; i < len(mList); i++ {
		N.Mul(N, mList[i])
	}
	return N
}

// 中国剩余定理
type CRT struct {
	N     *big.Int   //N=m1*m2*...*mn
	mList []*big.Int //两两互素，本身不一定是素数
	//Mi    []*big.Int //Mi=m1*..*mi-1,mi+1*..*mn
	//MiInv []*big.Int //Mi的逆 mod mi
	mimulmiinv []*big.Int
}

// 生成 n 个不同的素数
func genNPrime(n, nBits int) ([]*big.Int, error) {
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

// 创建一个CRT对象
func NewCRT(mList []*big.Int) *CRT {
	N := MulBigIntList(mList)
	Mi, MiInv := MiAndInvFrommList(mList)
	mimulmiinv := MiMulMiinv(Mi, MiInv, N)
	return &CRT{
		N:          N,
		mList:      mList,
		mimulmiinv: mimulmiinv,
	}
}
func (c *CRT) Set(mList []*big.Int) {
	N := MulBigIntList(mList)
	Mi, MiInv := MiAndInvFrommList(mList)
	mimulmiinv := MiMulMiinv(Mi, MiInv, N)
	c.N = N
	c.mimulmiinv = mimulmiinv
	c.mList = mList
}
func N2Zmi(x *big.Int, mlist []*big.Int) []*big.Int {
	n := len(mlist)
	res := make([]*big.Int, 0, n)
	for i := 0; i < n; i++ {
		a := new(big.Int).Mod(x, mlist[i])
		res = append(res, a)
	}
	return res
}

// 中国剩余定理，Zn-->Zmi
func (c *CRT) N2Zmi(x *big.Int) []*big.Int {
	return N2Zmi(x, c.mList)
}

// 中国剩余定理，Zpi-->Zn
func (c *CRT) Zmi2N(ai []*big.Int) *big.Int {
	n := len(ai)
	if n != len(c.mList) {
		panic("error parameter")
	}
	sum := big.NewInt(0)
	tmp := new(big.Int)
	for i := 0; i < n; i++ {
		//tmp := MulMod(c.Mi[i], c.MiInv[i], c.N)
		tmp.Mul(c.mimulmiinv[i], ai[i])
		tmp.Mod(tmp, c.N)
		sum.Add(sum, tmp)
		sum.Mod(sum, c.N)
	}
	return sum
}

func MulMod(x, y, P *big.Int) *big.Int {
	tmp := new(big.Int).Mul(x, y)
	tmp.Mod(tmp, P)
	return tmp
}

// 从模N的既约剩余系中随机选择一个元素
// reduced residue system ,RRS 既约剩余系
func (c *CRT) RandomFromZnReduced() *big.Int {
	ai := make([]*big.Int, 0, len(c.mList))
	for i := 0; i < len(c.mList); i++ {
		ai = append(ai, RandomFromZnNotZero(c.mList[i]))
	}
	z := c.Zmi2N(ai)
	return z
}
func RandomFromZnReduced(N *big.Int) *big.Int {
	count := 0
	for true {
		count++
		x := RandomFromZnNotZero(N)
		gcd := GcdEuclid(x, N)
		if gcd.Cmp(One) == 0 {
			//log.Printf("=== count:%+v\n", count)
			return x
		}
	}
	return nil
}

// 获取Mi以及Mi^{-1}
func (c *CRT) MiAndInvFrommList() ([]*big.Int, []*big.Int) {
	return MiAndInvFrommList(c.mList)

}

/******************************** crt *********************************/
// 获取Mi以及Mi^{-1}
func MiAndInvFrommList(pList []*big.Int) ([]*big.Int, []*big.Int) {
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
func MiMulMiinv(Mi []*big.Int, Miinv []*big.Int, N *big.Int) []*big.Int {
	mimulmiinv := make([]*big.Int, 0, len(Mi))
	for i := 0; i < len(Mi); i++ {
		tmp := new(big.Int).Mul(Mi[i], Miinv[i])
		tmp.Mod(tmp, N)
		mimulmiinv = append(mimulmiinv, tmp)
	}
	return mimulmiinv
}
func GenMiMulMiinv(mlist []*big.Int, N *big.Int) []*big.Int {
	Mi, Inv := MiAndInvFrommList(mlist)
	return MiMulMiinv(Mi, Inv, N)
}
func ZmiToN(ai []*big.Int, mimulmiinv []*big.Int, N *big.Int) *big.Int {
	sum := big.NewInt(0)
	tmp := new(big.Int)
	for i := 0; i < len(ai); i++ {
		tmp.Mul(ai[i], mimulmiinv[i])
		tmp.Mod(tmp, N)
		sum.Add(sum, tmp)
		sum.Mod(sum, N)
	}
	return sum
}
func NToZmi(x *big.Int, mlist []*big.Int) []*big.Int {
	n := len(mlist)
	ai := make([]*big.Int, 0, n)
	for i := 0; i < n; i++ {
		a := new(big.Int).Mod(x, mlist[i])
		ai = append(ai, a)
	}
	return ai
}
