package alg

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestRsa2_Gen(t *testing.T) {
	var rs Rsa2
	for i := 0; i < 10; i++ {
		pk, prik, err := rs.Gen(512)
		assert.NoError(t, err)
		tmp := new(big.Int).Mul(pk.e, prik.d)
		tmp.Mod(tmp, prik.faiN)
		fg := tmp.Cmp(One)
		assert.Equal(t, 0, fg)
	}
}
func TestRsa2_Enc(t *testing.T) {
	var r Rsa2
	pk, sk, err := r.Gen(512)
	assert.NoError(t, err)
	for i := 0; i < 100000; i++ {
		plain := RandomFromZnNotZero(pk.N)
		cip := r.Enc(pk, plain)
		//mpalin := r.Dec(sk, cip)
		mpalin := r.Dec2(sk, cip)
		assert.Equal(t, 0, plain.Cmp(mpalin))
	}
}
