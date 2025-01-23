package alg

import (
	"math/big"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestGcdEuclid(t *testing.T) {
	assert.New(t)
	a, b := big.NewInt(12), big.NewInt(18)
	ret := GcdEuclid(a, b)
	assert.Equal(t, big.NewInt(6), ret)
	assert.Equal(t, big.NewInt(4), GcdEuclid(big.NewInt(12), big.NewInt(20)))
}
func TestGenRandomBytes(t *testing.T) {
	nBytes := 19
	ret := GenRandomBytes(nBytes)
	assert.Equal(t, len(ret), nBytes)
	assert.Equal(t, len(GenRandomBytes(5)), 5)
	assert.Equal(t, len(GenRandomBytes(15)), 15)
	assert.Equal(t, len(GenRandomBytes(11)), 11)

}
func TestGenRandomBigInt(t *testing.T) {
	nbit := 9
	ret := GenRandomBigInt(nbit)
	assert.Equal(t, nbit, ret.BitLen())
}
func TestGenRandomFromZn(t *testing.T) {
	ret := GenRandomFromZn(big.NewInt(3))
	assert.Equal(t, ret.Cmp(big.NewInt(3)), -1)
	ret = GenRandomFromZn(big.NewInt(7))
	assert.Equal(t, ret.Cmp(big.NewInt(7)), -1)
	ret = GenRandomFromZn(big.NewInt(256))
	assert.Equal(t, ret.Cmp(big.NewInt(256)), -1)
	ret2 := new(big.Int).Div(big.NewInt(17), big.NewInt(4))
	assert.Equal(t, ret2.String(), "4")
}
func TestMillerRabin(t *testing.T) {
	ret := MillerRabin(big.NewInt(17), 40)
	assert.True(t, ret)
	ret = MillerRabin(big.NewInt(25), 40)
	assert.False(t, ret)
	ret = MillerRabin(big.NewInt(19), 40)
	assert.True(t, ret)
}

func TestGenPrime(t *testing.T) {
	N, err := GenPrime(381)
	assert.NoError(t, err)
	assert.Equal(t, N.BitLen(), 381)
}
func TestGetNPrime(t *testing.T) {
	res, err := GetNPrime(10, 128)
	assert.NoError(t, err)
	assert.Equal(t, len(res), 10)
}
func TestNewCRT(t *testing.T) {
	crt := NewCRT(5, 8)
	assert.NotNil(t, crt)
	assert.Equal(t, len(crt.pList), 5)
}
func TestNewCRT2(t *testing.T) {
	crt := NewCRT(5, 8)
	assert.NotNil(t, crt)
	crt2 := NewCRT2(crt.pList)
	assert.NotNil(t, crt2)
}

func TestCRT_N2Zpi(t *testing.T) {
	n := 2
	crt := NewCRT(n, 32)
	assert.NotNil(t, crt)
	x := GenRandomFromZn(crt.N)
	assert.Equal(t, x.Cmp(crt.N), -1)
	ai := crt.N2Zpi(x)
	//Mi, MiInv := crt.GetMiAndInvFromPList()
	x2 := crt.Zpi2N(ai)
	assert.Equal(t, x.String(), x2.String())
}
func TestGetMiAndInvFromPList(t *testing.T) {
	qi := []*big.Int{
		big.NewInt(3),
		big.NewInt(5),
		big.NewInt(7),
	}
	Mi, MiInv := GetMiAndInvFromPList(qi)
	assert.Equal(t, Mi[0].String(), "35")
	assert.Equal(t, Mi[1].String(), "21")
	assert.Equal(t, Mi[2].String(), "15")
	tmp := new(big.Int).Mul(Mi[0], MiInv[0])
	assert.Equal(t, tmp.Mod(tmp, qi[0]), One)
	tmp = new(big.Int).Mul(Mi[1], MiInv[1])
	assert.Equal(t, tmp.Mod(tmp, qi[1]), One)
	tmp = new(big.Int).Mul(Mi[2], MiInv[2])
	assert.Equal(t, tmp.Mod(tmp, qi[2]), One)
	N := big.NewInt(3 * 5 * 7)
	x := GenRandomFromZn(N)
	assert.Equal(t, x.Cmp(N), -1)
	crt := NewCRT2(qi)
	assert.NotNil(t, crt)
	ai := crt.N2Zpi(x)
	x2 := crt.Zpi2N(ai)
	assert.Equal(t, x.String(), x2.String())
}

func BenchmarkCRT_N2Zpi(b *testing.B) {
	n, nBits := 2, 32
	for i := 0; i < b.N; i++ {
		crt := NewCRT(n, nBits)
		x := GenRandomFromZn(crt.N)
		ai := crt.N2Zpi(x)
		crt.Zpi2N(ai)
	}
}

func BenchmarkGenPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenPrime(128)
	}
}
func BenchmarkGenRandomBigInt(b *testing.B) {
	nbit := 128
	for i := 0; i < b.N; i++ {
		GenRandomBigInt(nbit)
	}
}
func BenchmarkMillerRabin(b *testing.B) {
	N := GenRandomBigInt(128)
	for i := 0; i < b.N; i++ {
		MillerRabin(N, 40)
	}
}
func BenchmarkGenRandomFromZn(b *testing.B) {
	N := GenRandomBigInt(123)
	for i := 0; i < b.N; i++ {
		GenRandomFromZn(N)
	}
}
func BenchmarkNewCRT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewCRT(2, 256)
	}
}
func BenchmarkGetNPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetNPrime(2, 256)
	}
}
