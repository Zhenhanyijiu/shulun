package alg

import (
	"math/big"
)

// 素数模的雅可比符号
func Jacobi(a, p *big.Int) *big.Int {
	tmp := new(big.Int).Div(new(big.Int).Sub(p, One), Two)
	return new(big.Int).Exp(a, tmp, p)
}

// 合数模N=pq形式的雅可比符号
func JacobiNpq(a, p, q *big.Int) *big.Int {
	return new(big.Int).Mul(Jacobi(a, p), Jacobi(a, q))
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
	//log.Printf("GenPrimeList use time:%+v\n", time.Since(start))
	c := NewCRT(mlist)
	//log.Printf("NewCRT use time:%+v\n", time.Since(start))
	xpn := RandomQNRModPrime(c.mList[0])
	//log.Printf("RandomQNRModPrime 1 use time:%+v\n", time.Since(start))
	xqn := RandomQNRModPrime(c.mList[1])
	//log.Printf("RandomQNRModPrime 2 use time:%+v\n", time.Since(start))
	z := c.Zmi2N([]*big.Int{xpn, xqn})
	//log.Printf("Zmi2N use time:%+v\n", time.Since(start))
	return &PubKey{c.N, z}, &PriKey{c, z}, nil
}

// 加密算法，明文只有 0 or 1
func (g *GoldwasserMicali) Enc(plaintext bool, N, z *big.Int) *big.Int {
	if plaintext { //加密 1

	} else { //加密 0

	}
	return nil
}
