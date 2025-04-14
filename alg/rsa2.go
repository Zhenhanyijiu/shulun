package alg

import (
	"math/big"
)

type Rsa2 struct {
}
type PubKey2 struct {
	N, e *big.Int
}
type PriKey2 struct {
	N, d      *big.Int
	pq        []*big.Int
	faiN      *big.Int
	dmodpq_1  []*big.Int
	mimulminv []*big.Int
}

func (r *Rsa2) Gen(nBit int) (*PubKey2, *PriKey2, error) {
	//教科书式的rsa,e固定3,p,q!=1mod3,才能保证e与fai(N)互素
	e, primeLen := big.NewInt(3), nBit/2
	plist, err := GenPrimeList(2, primeLen)
	if err != nil {
		return nil, nil, err
	}
	tmp := new(big.Int)
	count := 0
	for tmp.Mod(plist[0], e).Cmp(One) == 0 || tmp.Mod(plist[1], e).Cmp(One) == 0 {
		count++
		plist, err = GenPrimeList(2, primeLen)
		if err != nil {
			return nil, nil, err
		}
	}
	//log.Printf("=== count:%+v ,plist:%+v\n", count, plist)
	//log.Printf("p bitlen:%+v,q bitlen:%+v", plist[0].BitLen(), plist[1].BitLen())
	N := new(big.Int).Mul(plist[0], plist[1])
	pq_1 := []*big.Int{new(big.Int).Sub(plist[0], One), new(big.Int).Sub(plist[1], One)}
	faiN := new(big.Int).Mul(pq_1[0], pq_1[1])
	d := new(big.Int).ModInverse(e, faiN)
	dmodpq_1 := []*big.Int{pq_1[0].Mod(d, pq_1[0]), pq_1[1].Mod(d, pq_1[1])}
	mimulminv := GenMiMulMiinv(plist, N)
	sk := &PriKey2{N: N, d: d, pq: plist,
		faiN: faiN, dmodpq_1: dmodpq_1,
		mimulminv: mimulminv,
	}
	return &PubKey2{N: N, e: e}, sk, nil
}
func (r *Rsa2) Enc(key2 *PubKey2, plain *big.Int) *big.Int {
	//c=m^e
	cipher := new(big.Int).Exp(plain, key2.e, key2.N)
	return cipher
}
func (*Rsa2) Dec(key2 *PriKey2, cipher *big.Int) *big.Int {
	//m=c^d
	plain := new(big.Int).Exp(cipher, key2.d, key2.N)
	return plain
}

func (*Rsa2) Dec2(key2 *PriKey2, cipher *big.Int) *big.Int {
	//m=c^d
	ap := new(big.Int).Exp(cipher, key2.dmodpq_1[0], key2.pq[0])
	aq := new(big.Int).Exp(cipher, key2.dmodpq_1[1], key2.pq[1])
	plain := ZmiToN([]*big.Int{ap, aq}, key2.mimulminv, key2.N)
	return plain
}
