package alg

import (
	assert "github.com/stretchr/testify/assert"
	"log"
	"math/big"
	"math/rand"
	"testing"
)

/*
go test -run=. -v
go test -bench .
运行基准测试：go test -bench=
进阶参数
-benchtime t 指定基准测试的运行时间。例如，运行 5 秒：
go test -bench=. -benchtime=5s

-count n 多次运行基准测试以提高准确性：
go test -bench=. -count=3

-cpu n 设置使用的 CPU 核数：
go test -bench=. -cpu=1,2,4

-benchmem 查看内存分配情况：
go test -bench=. -benchmem
注意事项
ResetTimer: 在耗时准备工作后重置计时器。
StopTimer & StartTimer: 在每次函数执行前暂停和继续计时。
go test -v -benchmem -count=1 -cpu=1 -bench=BenchmarkGcdEuclid
*/
func genIntListTest(n int) []int64 {
	ret := make([]int64, 0)
	for i := 0; i < n; i++ {
		ret = append(ret, int64(rand.Uint32()%256))
	}
	return ret
}
func TestGcdEuclid(t *testing.T) {
	a, b := big.NewInt(12), big.NewInt(18)
	ret := GcdEuclid(a, b)
	assert.Equal(t, big.NewInt(6), ret)
	assert.Equal(t, big.NewInt(4), GcdEuclid(big.NewInt(12), big.NewInt(20)))
}
func BenchmarkGcdEuclid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GcdEuclid(big.NewInt(int64(976653400)), big.NewInt(int64(224723750)))
	}
}

func TestRandomBytes(t *testing.T) {
	nBytes := 19
	assert.Equal(t, len(RandomBytes(nBytes)), nBytes)
	assert.Equal(t, len(RandomBytes(5)), 5)
	assert.Equal(t, len(RandomBytes(15)), 15)
	assert.Equal(t, len(RandomBytes(11)), 11)

}
func BenchmarkRandomBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomBytes(128)
	}
}
func TestRandomBigInt(t *testing.T) {
	rnd := func() int {
		a := rand.Uint32()
		a = a % 8192
		return int(a) + 1
	}
	for i := 0; i < 100000; i++ {
		nbit := rnd()
		ret := RandomBigInt(nbit)
		assert.Equal(t, nbit, ret.BitLen())
	}
}
func BenchmarkRandomBigInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomBigInt(4096)
	}
}

// 包括 0
func TestRandomFromZn(t *testing.T) {
	n := big.NewInt(100)
	count := 0
	for i := 0; i < 10000*10; i++ {
		r := RandomFromZn(n)
		assert.Equal(t, 1, r.Cmp(NegOne))
		assert.Equal(t, -1, r.Cmp(n))
		count++
	}
	log.Printf("count:%+v\n", count)
}
func BenchmarkRandomFromZn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandomFromZn(big.NewInt(int64(1 << 3)))
		//_ = GenRandomFromZnNotZero2(big.NewInt(int64(1 << 3)))
	}
}

// 不包括 0
func TestRandomFromZnNotZero(t *testing.T) {
	n := big.NewInt(100)
	count := 0
	for i := 0; i < 10000*10; i++ {
		r := RandomFromZnNotZero(n)
		assert.Equal(t, 1, r.Cmp(Zero))
		assert.Equal(t, -1, r.Cmp(n))
		count++
	}
	log.Printf("count:%+v\n", count)
}
func BenchmarkRandomFromZnNotZero(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandomFromZnNotZero(big.NewInt(int64(1 << 3)))
		//_ = GenRandomFromZnNotZero2(big.NewInt(int64(1 << 3)))
	}
}
func TestNTods(t *testing.T) {
	count := 0
	for i := 0; i < 10*10000; i++ {
		maxn := rand.Int63()
		if maxn > 0 && maxn%2 == 0 {
			maxn2 := big.NewInt(maxn)
			d, s := NTods(maxn2)
			assert.Equal(t, 0, new(big.Int).Mod(d, Two).Cmp(One))
			assert.Less(t, 0, s)
			assert.Greater(t, s, 0)
			tmp := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(s)), new(big.Int).Add(maxn2, One))
			fg := tmp.Mul(tmp, d).Cmp(maxn2)
			assert.Equal(t, 0, fg)
			count++
		} else {
			i--
		}
	}
	log.Printf("count:%+v\n", count)
}
func BenchmarkNTods(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NTods(big.NewInt(100000000))
	}
}
func TestMillerRabin(t *testing.T) {
	P, ok := new(big.Int).SetString("186151550718860665121779391551741013053", 10)
	assert.True(t, ok)
	NotP, ok := new(big.Int).SetString("186151550718860665121779391551741013054", 10)
	assert.True(t, ok)
	for i := 0; i < 10*100; i++ {
		ret := MillerRabin(big.NewInt(25), 40)
		assert.False(t, ret)
		ret2 := MillerRabin(big.NewInt(9973), 40)
		assert.True(t, ret2)
		ret3 := MillerRabin(P, 40)
		assert.True(t, ret3)
		ret4 := MillerRabin(NotP, 40)
		assert.False(t, ret4)
	}
}

func BenchmarkMillerRabin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MillerRabin(BigPrime, 40)
		//MillerRabin(big.NewInt(41), 40)
	}
}

func TestGenPrime(t *testing.T) {
	for i := 0; i < 100; i++ {
		nbit := int(rand.Uint32()%512 + 5)
		N, err := GenPrime(nbit)
		assert.NoError(t, err)
		assert.Equal(t, N.BitLen(), nbit)
	}
}
func BenchmarkGenPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenPrime(256)
	}
}
func TestGenPrimeList(t *testing.T) {
	for i := 0; i < 20; i++ {
		plist, err := GenPrimeList(3, 256)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(plist))
		for i2 := 0; i2 < len(plist); i2++ {
			ok := MillerRabin(plist[i2], 40)
			assert.True(t, ok)
		}
	}
}
func BenchmarkGenPrimeList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenPrimeList(2, 256)
	}
}

func TestMulBigIntList(t *testing.T) {
	mlist := []*big.Int{big.NewInt(3), big.NewInt(4), big.NewInt(5)}
	N := MulBigIntList(mlist)
	assert.Equal(t, 0, N.Cmp(big.NewInt(3*4*5)))
	for i := 0; i < 128; i++ {
		intlist := genIntListTest(5)
		mlist := make([]*big.Int, 0)
		s := int64(1)
		for i2 := 0; i2 < len(intlist); i2++ {
			mlist = append(mlist, big.NewInt(intlist[i2]))
			s = s * intlist[i2]
		}
		N := MulBigIntList(mlist)
		assert.Equal(t, 0, N.Cmp(big.NewInt(s)))
	}
}

func TestCRT_N2Zmi(t *testing.T) {
	for i := 0; i < 1000; i++ {
		n := 3
		mlist, err := GenPrimeList(n, 8)
		assert.NoError(t, err)
		crt := NewCRT(mlist)
		assert.NotNil(t, crt)
		x := RandomFromZnNotZero(crt.N)
		assert.Equal(t, -1, x.Cmp(crt.N))
		assert.Equal(t, 1, x.Cmp(Zero))
		ai := crt.N2Zmi(x)
		x2 := crt.Zmi2N(ai)
		assert.Equal(t, 0, x2.Cmp(x))
	}
}
func BenchmarkCRT_N2Zmi(b *testing.B) {
	mlist, _ := GenPrimeList(3, 32)
	crt := NewCRT(mlist)
	for i := 0; i < b.N; i++ {
		_ = crt.N2Zmi(big.NewInt(10000))
	}
}
func TestCRT_Zmi2N(t *testing.T) {
	mlist, err := GenPrimeList(3, 32)
	assert.NoError(t, err)
	c := NewCRT(mlist)
	for i := 0; i < 1000; i++ {
		xn := RandomFromZnNotZero(c.N)
		ai := c.N2Zmi(xn)
		xn2 := c.Zmi2N(ai)
		assert.Equal(t, 0, xn.Cmp(xn2))
	}
}
func BenchmarkCRT_Zmi2N(b *testing.B) {
	mlist, _ := GenPrimeList(3, 32)
	c := NewCRT(mlist)
	xn := RandomFromZnNotZero(c.N)
	ai := c.N2Zmi(xn)
	for i := 0; i < b.N; i++ {
		_ = c.Zmi2N(ai)
	}
}
func TestCRT_RandomFromZnReduced(t *testing.T) {
	mlist, err := GenPrimeList(3, 32)
	assert.NoError(t, err)
	c := NewCRT(mlist)
	for i := 0; i < 1000; i++ {
		x := c.RandomFromZnReduced()
		gcd := GcdEuclid(x, c.N)
		assert.Equal(t, 0, gcd.Cmp(One))
	}

}
func BenchmarkCRT_RandomFromZnReduced(b *testing.B) {
	mlist, _ := GenPrimeList(3, 128)
	c := NewCRT(mlist)
	for i := 0; i < b.N; i++ {
		_ = c.RandomFromZnReduced()
	}
}
func TestRandomFromZnReduced(t *testing.T) {
	mlist, err := GenPrimeList(2, 128)
	assert.NoError(t, err)
	c := NewCRT(mlist)
	for i := 0; i < 1; i++ {
		x := RandomFromZnReduced(c.N)
		gcd := GcdEuclid(x, c.N)
		assert.Equal(t, 0, gcd.Cmp(One))
	}
}
func BenchmarkRandomFromZnReduced(b *testing.B) {
	mlist, _ := GenPrimeList(3, 128)
	c := NewCRT(mlist)
	for i := 0; i < b.N; i++ {
		_ = RandomFromZnReduced(c.N)
	}
}
func TestCRT_MiAndInvFrommList(t *testing.T) {
	for i := 0; i < 1000; i++ {
		mlist, err := GenPrimeList(3, 32)
		assert.NoError(t, err)
		Mi, Inv := MiAndInvFrommList(mlist)
		for i2 := 0; i2 < len(mlist); i2++ {
			tmp := new(big.Int).Mul(Mi[i2], Inv[i2])
			tmp.Mod(tmp, mlist[i2])
			assert.Equal(t, 0, tmp.Cmp(One))
		}
	}
}

func BenchmarkCRT_MiAndInvFrommList(b *testing.B) {
	mlist, _ := GenPrimeList(3, 128)
	for i := 0; i < b.N; i++ {
		_, _ = MiAndInvFrommList(mlist)
	}
}
