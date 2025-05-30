package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"shulun/alg"
	"time"
)

func main() {
	log.Printf("------------- main ---------------\n")
	//test_genprime(20)
	//test_rand_prime(20)
	//test_generater(1)
	//test_ecc_pairing()
	//test_prime_field()
	//mp := QR(big.NewInt(5))
	//fmt.Printf("mp:%+v\n", mp)
	EccTestTest()
}
func getCoffFp(cof []int64, p int64) ([]*big.Int, *big.Int) {
	coff := make([]*big.Int, 0, len(cof))
	for i := 0; i < len(cof); i++ {
		coff = append(coff, big.NewInt(cof[i]))
	}
	return coff, big.NewInt(p)
}
func EccTestTest() {
	//coff, Fp := getCoffFp([]int64{0, 1, 0, 1}, 5)
	//coff, Fp := getCoffFp([]int64{3, 2, 0, 1}, 7)
	coff, Fp := getCoffFp([]int64{3, 1, 0, 1}, 11)
	//coff, Fp := getCoffFp([]int64{2, 7, 0, 1}, 11)
	//coff, Fp := getCoffFp([]int64{7, 13, 0, 1}, 23)
	ecc := NewEccTest(coff, Fp)
	ecc.GenAllPoint()
	ecc.PrintInfo()
}

type EccTest struct {
	coff []*big.Int
	Fp   *big.Int
	qr   map[string]string
	all  []*Point
}
type Point struct {
	x, y *big.Int
}

func NewEccTest(coff []*big.Int, p *big.Int) *EccTest {
	qr := QR(p)
	return &EccTest{
		coff: coff, Fp: p, qr: qr}
}
func (e *EccTest) GenAllPoint() {
	count := e.Fp.Int64()
	for i := int64(0); i < count; i++ {
		x := big.NewInt(i)
		yqr := poly(e.coff, x, e.Fp)
		key := yqr.String()
		val, ok := e.qr[key]
		if ok {
			y, fg := new(big.Int).SetString(val, 10)
			if fg {
				if y.Cmp(big.NewInt(0)) == 0 {
					e.all = append(e.all, &Point{x: x, y: y})
					continue
				}
				e.all = append(e.all, &Point{x: x, y: y})
				y1 := new(big.Int).Neg(y)
				y1.Mod(y1, e.Fp)
				e.all = append(e.all, &Point{x: x, y: y1})
			} else {
				panic("error setstring!")
			}
		}

	}
	e.all = append(e.all, &Point{})
}
func (e *EccTest) PrintInfo() {
	fmt.Printf("coff:%+v\n", e.coff)
	fmt.Printf("Fp:%+v\n", e.Fp)
	fmt.Printf("Qr:%+v\n", e.qr)

	fmt.Printf("all point:%v\n", len(e.all))
	for i := 0; i < len(e.all); i++ {
		fmt.Printf("(%+v,%+v),", e.all[i].x, e.all[i].y)
	}
	fmt.Println()
}
func test_prime_field() {
	n := 19
	zp_star := []int64{}
	for i := 1; i < n; i++ {
		zp_star = append(zp_star, int64(i))
	}
	fmt.Printf("zp_star:%+v\n", zp_star)
	m := big.NewInt(int64(n))
	for i := 0; i < len(zp_star); i++ {
		tmp := big.NewInt(zp_star[i])
		arr := make([]int64, 0, len(zp_star))
		for i2 := 0; i2 < len(zp_star); i2++ {
			tmp2 := new(big.Int).Exp(tmp, big.NewInt(int64(i2)), m)
			arr = append(arr, tmp2.Int64())
		}
		fmt.Printf("%+v,%+v\n", zp_star[i], arr)
	}
}
func test_ecc_pairing() {
	n := 11
	q := big.NewInt(int64(n))
	coff := []*big.Int{big.NewInt(2), big.NewInt(7), big.NewInt(0), big.NewInt(1)}
	ret := poly(coff, big.NewInt(1), q)
	fmt.Printf("ret:%+v\n", ret)
	ret = poly(coff, big.NewInt(2), q)
	fmt.Printf("ret:%+v\n", ret)
	QR := []*big.Int{}
	for i := 0; i < n/2+1; i++ {
		tmp := new(big.Int).Exp(big.NewInt(int64(i)), big.NewInt(2), q)
		fmt.Printf("i:%+v,tmp:%+v\n", i, tmp)
		QR = append(QR, tmp)
	}
	fmt.Printf("QR:%+v\n", QR)
	for i := 0; i < n; i++ {
		ret := poly(coff, big.NewInt(int64(i)), q)
		fmt.Printf("i:%+v,ret:%+v\n", i, ret)
	}
}
func poly(coff []*big.Int, x, q *big.Int) *big.Int {
	n := len(coff)
	tmp := new(big.Int).Set(coff[n-1])
	for i := n - 1; i > 0; i-- {
		//ai*x+ai_1
		tmp.Mul(tmp, x).Mod(tmp, q)
		tmp.Add(tmp, coff[i-1]).Mod(tmp, q)
	}
	return tmp
}
func QR(prime *big.Int) map[string]string {
	mp := make(map[string]string)
	zero := big.NewInt(0).String()
	mp[zero] = zero
	count := prime.Int64() - 1
	x := big.NewInt(0)
	for i := int64(1); i <= count/2; i++ {
		x.SetInt64(i)
		tmp := big.NewInt(i)
		tmp.Exp(tmp, big.NewInt(2), prime)
		mp[tmp.String()] = x.String()
	}

	return mp
}

/*
def mod_pow(a: dec.Decimal, n: dec.Decimal, p: dec.Decimal):

	# 将指数n转换为二进制数组
	n_binary = bin(n)[2:]
	result, base = 1, a % p
	# 从最低位开始遍历n的二进制表示
	for digit in reversed(n_binary):
	    if digit == '1':
	        result = (result * base) % p
	    base = (base * base) % p
	return result
*/
func mod_pow(a, n, p int) int {
	ret := fmt.Sprintf("%b", n)
	//log.Printf("ret:%+v\n", ret)
	result, base := 1, a%p
	for i := len(ret) - 1; i >= 0; i-- {
		if ret[i] == '1' {
			result = (result * base) % p
		}
		base = (base * base) % p
	}
	return result
}
func test_generater(cyc int) {
	a := 2
	p, q := 5, 7
	N := p * q
	fai := (p - 1) * (q - 1)
	for i := 0; i < fai; i++ {
		fmt.Printf("%v^%v=%+v\n", a, i, mod_pow(a, i, N))
	}
	N1 := big.NewInt(int64(N))
	faiInv := new(big.Int).ModInverse(big.NewInt(int64(fai)), N1)
	fmt.Printf("fai inv:%+v\n", faiInv)

}
func test_genprime(cyc int) {
	log.Printf("--------- test_genprime ----------")
	for i := 0; i < cyc; i++ {
		start := time.Now()
		_, err := alg.GenPrime(512)
		log.Printf("use time:%+v,error:%+v\n", time.Since(start), err)
	}
}

func test_rand_prime(cyc int) {
	log.Printf("--------- test_rand_prime ----------")
	for i := 0; i < cyc; i++ {
		start := time.Now()
		_, err := rand.Prime(rand.Reader, 512)
		log.Printf("use time:%+v,error:%+v\n", time.Since(start), err)
	}
}
