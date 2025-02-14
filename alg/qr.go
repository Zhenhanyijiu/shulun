package alg

import (
	"math/big"
)

// 素数模的雅可比符号
func Jacobi(a, p *big.Int) *big.Int {
	tmp := new(big.Int).Sub(p, One)
	tmp.Div(tmp, Two)
	return tmp.Exp(a, tmp, p)
}

// 合数模N=pq形式的雅可比符号
func JacobiNpq(a, p, q *big.Int) *big.Int {
	return Jacobi(a, p).Mul(Jacobi(a, p), Jacobi(a, q))
}

/*
P是一个素数，判断Zp中的元素是否是二次剩余
判定素数模二次剩余
*/
func IsQRModPrime(a, P *big.Int) bool {
	if Jacobi(a, P).Cmp(One) == 0 {
		return true
	}
	return false
}

// 判定某数是否是合数模的二次剩余,p,q都是素数
func IsQRModComposite(x, p, q *big.Int) bool {
	if Jacobi(x, p).Cmp(One) == 0 && Jacobi(x, q).Cmp(One) == 0 {
		return true
	}
	return false
}

// 素数模下生成一个二次非剩余
func RandomQNRModPrime(p *big.Int) *big.Int {
	count := 0
	for true {
		count++
		tmp := RandomFromZnNotZero(p)
		t2 := Jacobi(tmp, p)
		if t2.Cmp(One) != 0 {
			//log.Printf("count:%+v\n", count)
			return tmp
		}
	}
	panic("not found")
}

type GoldwasserMicali struct {
	//c *CRT
}
type PriKey struct {
	c *CRT
	z *big.Int
}
type PubKey struct {
	N, z *big.Int
}

func NewGoldwasserMicali() *GoldwasserMicali {
	return &GoldwasserMicali{}
}

func (g *GoldwasserMicali) Gen(primeBitLen int) (*PubKey, *PriKey, error) {
	//start := time.Now()
	mlist, err := GenPrimeList(2, primeBitLen)
	if err != nil {
		return nil, nil, err
	}
	c := NewCRT(mlist)
	xpn := RandomQNRModPrime(c.mList[0])
	xqn := RandomQNRModPrime(c.mList[1])
	z := c.Zmi2N([]*big.Int{xpn, xqn})
	return &PubKey{c.N, z}, &PriKey{c, z}, nil
}

// 加密算法，明文只有 0 or 1
// 二次剩余 对应明文 0；二次非剩余 对应明文 1；
func (p *PubKey) Enc(plaintext bool) *big.Int {
	x := RandomFromZnReduced(p.N)
	c := x.Exp(x, Two, p.N)
	//x.Mod(x, x)
	//c := x.Mod(x, p.N)
	if plaintext { //加密 1
		c.Mul(p.z, c)
	} //加密 0
	return c
}
func (k *PriKey) Dec(c *big.Int) bool {
	fg := JacobiNpq(c, k.c.mList[0], k.c.mList[1])
	plain := true
	if fg.Cmp(One) == 0 {
		plain = false
	}
	return plain

}
