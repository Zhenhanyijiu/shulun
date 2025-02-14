package main

import (
	"crypto/rand"
	"log"
	"shulun/alg"
	"time"
)

func main() {
	log.Printf("------------- main ---------------\n")
	test_genprime(20)
	test_rand_prime(20)
}
func test_genprime(cyc int) {
	log.Printf("--------- test_genprime ----------")
	for i := 0; i < cyc; i++ {
		start := time.Now()
		_, err := alg.GenPrime(3072)
		log.Printf("use time:%+v,error:%+v\n", time.Since(start), err)
	}
}

func test_rand_prime(cyc int) {
	log.Printf("--------- test_rand_prime ----------")
	for i := 0; i < cyc; i++ {
		start := time.Now()
		_, err := rand.Prime(rand.Reader, 3072)
		log.Printf("use time:%+v,error:%+v\n", time.Since(start), err)
	}
}
