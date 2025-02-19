package alg

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

// P是奇数素数，求出二次剩余和二次非剩余,用于测试
func qrAndQNRTest(p int) ([]int, []int) {
	tmp := make([]bool, p)
	for i := 1; i < p; i++ {
		tmp[i] = false
	}
	for i := 1; i < p; i++ {
		tmp[i*i%p] = true
	}
	qr, qnr := make([]int, 0), make([]int, 0)
	for i := 1; i < p; i++ {
		if tmp[i] {
			qr = append(qr, i)
		} else {
			qnr = append(qnr, i)
		}
	}
	//log.Printf("qr :%+v\n", qr)
	//log.Printf("qnr:%+v\n", qnr)
	return qr, qnr
}

func TestJacobi(t *testing.T) {
	for i := 0; i < 10; i++ {
		p, err := GenPrime(14)
		assert.NoError(t, err)
		//log.Printf("==> p:%+v\n", p)
		qr, qnr := qrAndQNRTest(int(p.Uint64()))
		assert.Equal(t, len(qr), len(qnr))
		for k := 0; k < len(qr); k++ {
			isqr := Jacobi(big.NewInt(int64(qr[k])), p)
			assert.Equal(t, 0, isqr.Cmp(One))
			notqr := Jacobi(big.NewInt(int64(qnr[k])), p)
			assert.Equal(t, 0, notqr.Cmp(new(big.Int).Mod(NegOne, p)))
		}
	}
}
func BenchmarkJacobi(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Jacobi(RandomFromZnNotZero(BigPrime), BigPrime)
	}
}
func TestJacobiNpq(t *testing.T) {
	for i := 0; i < 10; i++ {
		mlist, err := GenPrimeList(2, 32)
		assert.NoError(t, err)
		//mlist := []*big.Int{big.NewInt(5), big.NewInt(7)}
		c := NewCRT(mlist)
		a := c.RandomFromZnReduced()
		fg := JacobiNpq(a, mlist[0], mlist[1])
		//log.Printf(">>>>>>>>> a:%+v,fg:%+v,p:%+v,q:%+v\n", a, fg, mlist[0], mlist[1])
		p_1 := new(big.Int).Mod(NegOne, mlist[0])
		q_1 := new(big.Int).Mod(NegOne, mlist[1])
		faiN := new(big.Int).Mul(p_1, q_1)
		if fg.Cmp(One) == 0 || fg.Cmp(p_1) == 0 || fg.Cmp(q_1) == 0 || fg.Cmp(faiN) == 0 {
		} else {
			assert.NoError(t, errors.New("error test"))
		}
	}
}
func BenchmarkJacobiNpq(b *testing.B) {
	mlist, _ := GenPrimeList(2, 128)
	c := NewCRT(mlist)
	a := c.RandomFromZnReduced()
	for i := 0; i < b.N; i++ {
		JacobiNpq(a, mlist[0], mlist[1])
	}
}
func TestRandomQNRModPrime(t *testing.T) {
	q := big.NewInt(11) //4229
	xqnr := RandomQNRModPrime(q)
	assert.False(t, IsQRModPrime(xqnr, q))
	for i := 0; i < 100; i++ {
		p, err := GenPrime(128)
		assert.NoError(t, err)
		Xqnr := RandomQNRModPrime(p)
		assert.False(t, IsQRModPrime(Xqnr, p))
	}
}
func BenchmarkRandomQNRModPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomQNRModPrime(BigPrime)
	}
}
func TestRandomQRmodPrime(t *testing.T) {
	for i := 0; i < 100; i++ {
		p, err := GenPrime(128)
		assert.NoError(t, err)
		xqr := RandomQRmodPrime(p)
		assert.True(t, IsQRModPrime(xqr, p))
	}
}
func TestGoldwasserMicali_Gen(t *testing.T) {
	for i := 0; i < 100; i++ {
		pk, sk, err := new(GoldwasserMicali).Gen(128)
		assert.NoError(t, err)
		assert.NotNil(t, pk)
		assert.NotNil(t, sk)
	}
}

func BenchmarkGoldwasserMicali_Gen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		new(GoldwasserMicali).Gen(512)
	}
}
func TestPubKey_Enc(t *testing.T) {
	pk, sk, err := new(GoldwasserMicali).Gen(32)
	assert.NoError(t, err)
	for i := 0; i < 10000; i++ {
		plain := ((i % 2) != 0)
		c := pk.Enc(plain)
		plain2 := sk.Dec(c)
		assert.Equal(t, plain, plain2)
	}
}
func BenchmarkPubKey_Enc(b *testing.B) {
	pk, _, _ := new(GoldwasserMicali).Gen(128)
	for i := 0; i < b.N; i++ {
		pk.Enc(true)
		pk.Enc(false)
	}
}
func genp1p3(nBit int) (*big.Int, *big.Int, error) {
	var p1, p3 *big.Int = nil, nil
	for p1 == nil || p3 == nil {
		x, err := GenPrime(nBit)
		if err != nil {
			return nil, nil, err
		}
		tmp := new(big.Int).Mod(x, Four)
		if tmp.Cmp(Three) == 0 {
			if p3 == nil {
				p3 = x
			}
		} else {
			if p1 == nil {
				p1 = x
			}
		}

	}
	return p1, p3, nil
}
func TestSqrtModPrime(t *testing.T) {
	for i := 0; i < 128; i++ {
		p1, p3, err := genp1p3(128)
		assert.NoError(t, err)
		//log.Printf("p1:%+v,p3:%+v\n", p1, p3)
		assert.Equal(t, 0, new(big.Int).Mod(p1, Four).Cmp(One))
		assert.Equal(t, 0, new(big.Int).Mod(p3, Four).Cmp(Three))
		xqr3 := RandomQRmodPrime(p3)
		xsqrt3 := SqrtModPrime(xqr3, p3)
		//log.Printf("xqr3:%+v,xsqrt3:%+v\n", xqr3, xsqrt3)
		assert.Equal(t, 0, xsqrt3.Exp(xsqrt3, Two, p3).Cmp(xqr3))

		xqr1 := RandomQRmodPrime(p1)
		xsqrt1 := SqrtModPrime(xqr1, p1)
		//log.Printf("xqr1:%+v,xsqrt1:%+v\n", xqr1, xsqrt1)
		assert.Equal(t, 0, xsqrt1.Exp(xsqrt1, Two, p1).Cmp(xqr1))
	}
}
func BenchmarkSqrtModPrime(b *testing.B) {
	p1, p3, _ := genp1p3(128)
	xqr1 := RandomQRmodPrime(p1)
	xqr3 := RandomQRmodPrime(p3)
	for i := 0; i < b.N; i++ {
		SqrtModPrime(xqr1, p1)
		SqrtModPrime(xqr3, p3)
	}
}
