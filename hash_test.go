package gopool

import (
	"fmt"
	"github.com/geedchin/go-task-pool/util"
	"math/rand"
	"sort"
	"testing"
)

type tmpX struct {
	id    int
	times int
}
type tmpXs []tmpX

func (t tmpXs) Len() int {
	return len(t)
}

func (t tmpXs) Less(i, j int) bool {
	if t[i].times == t[j].times {
		return t[i].id < t[j].id
	}
	return t[i].times > t[j].times
}

func (t tmpXs) Swap(i, j int) {
	tmp := t[i]
	t[i] = t[j]
	t[j] = tmp
}

var (
	primes = []int{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47,
		53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113,
		127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197,
		199, 211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281,
		283, 293, 307, 311, 313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379,
		383, 389, 397, 401, 409, 419, 421, 431, 433, 439, 443, 449, 457, 461, 463,
		467, 479, 487, 491, 499, 503, 509, 521, 523, 541, 547, 557, 563, 569, 571,
		577, 587, 593, 599, 601, 607, 613, 617, 619, 631, 641, 643, 647, 653, 659,
		661, 673, 677, 683, 691, 701, 709, 719, 727, 733, 739, 743, 751, 757, 761,
		769, 773, 787, 797, 809, 811, 821, 823, 827, 829, 839, 853, 857, 859, 863,
		877, 881, 883, 887, 907, 911, 919, 929, 937, 941, 947, 953, 967, 971, 977,
		983, 991, 997,
	}
)

func TestNextPrim(t *testing.T) {

	idx := 0
	for i := 0; i <= 997 && idx < len(primes); i++ {
		if i > primes[idx] {
			idx++
		}
		if idx >= len(primes) {
			break
		}
		if tmp := nextPrim(i); tmp != primes[idx] {
			t.Errorf("i:%d, nextPrime want %d,but get %d",
				i, primes[idx], tmp)
		}
	}
}

func TestIsPrime(t *testing.T) {
	pm := make(map[int]bool, len(primes))
	for _, prime := range primes {
		pm[prime] = true
	}
	for i := 0; i < 1000; i++ {
		result := isPrime(i)
		want := pm[i]
		if result == want {
			continue
		}
		t.Errorf("testIsPrime want %v,but get %v", want, result)
	}
}

func TestHash(t *testing.T) {

	if tmp := hash(""); tmp != 0 {
		t.Errorf(`input "",want 0,but get %d`, tmp)
	}

	mp := make(map[int]int)
	size := 53
	for i := 0; i < 1000; i++ {
		result := hash(util.MD5([]byte(util.GenRandomDigitLowerLetter(rand.Int() % 51))))
		result %= size
		mp[result] += 1
	}
	arr := make(tmpXs, 0, len(mp))
	for k, v := range mp {
		arr = append(arr, tmpX{id: k, times: v})
	}
	sort.Sort(arr)
	for i := 0; i < arr.Len(); i++ {
		fmt.Println(arr[i].id, arr[i].times)
	}
	fmt.Println()
	fmt.Println(arr.Len())
}
