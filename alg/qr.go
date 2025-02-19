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

// 素数模 P 中随机生成一个二次剩余
func RandomQRmodPrime(p *big.Int) *big.Int {
	tmp := RandomFromZnNotZero(p)
	//log.Printf("tmpx:%+v\n", tmp)
	return tmp.Exp(tmp, Two, p)
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

// 找一个素数模 p 的二次非剩余
// 求素数模的平方根,p是一个奇素数，a是一个二次剩余
func SqrtModPrime(a, p *big.Int) *big.Int {
	r := new(big.Int).Mod(p, Four)
	//log.Printf(">>>>>>>>p mod4:%+v\n", r)
	//情况1：p%4==3
	if r.Cmp(Three) == 0 {
		r = r.Add(p, One)
		r.Div(r, Four)
		//log.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>..>.....")
		return r.Exp(a, r, p)
	}
	//log.Printf("=========== p mod 4== 1 here")
	b := RandomQNRModPrime(p)
	r.Sub(p, One)
	r.Div(r, Two) //(p-1)/2
	pMinus1Div2 := new(big.Int).Set(r)
	_, l := NTods(r)
	r1 := big.NewInt(0)
	tmp, tmp1 := new(big.Int), new(big.Int)
	for i := l; i >= 1; i-- {
		r.Div(r, Two)
		r1.Div(r1, Two)
		tmp.Exp(a, r, p)
		tmp1.Exp(b, r1, p)
		tmp.Mul(tmp, tmp1)
		tmp.Mod(tmp, p)
		if tmp.Cmp(One) != 0 {
			r1.Add(r1, pMinus1Div2)
		}
	}
	tmp.Div(tmp.Add(r, One), Two)
	tmp.Exp(a, tmp, p)
	tmp1.Div(r1, Two)
	tmp1.Exp(b, tmp1, p)
	tmp.Mul(tmp, tmp1)
	return tmp.Mod(tmp, p)
}

// 求合数模 N=pq 的平方根，利用素数模和中国剩余定理很容易得出
func SqrtModN(a, p, q *big.Int) *big.Int {
	return nil
}
