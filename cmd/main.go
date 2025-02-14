package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"shulun/alg"
	"time"
)

func main() {
	log.Printf("------------- main ---------------\n")
	//test_genprime(20)
	//test_rand_prime(20)
	test_generater(1)
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
