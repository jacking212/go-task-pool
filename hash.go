package gopool

func hash(s string) int {
	if s == "" {
		return 0
	}
	id := 1

	sb := []byte(s)
	for _, b := range sb {
		id += int(b)
	}
	return id
}

func nextPrim(k int) int {
	if k < 3 {
		return 2
	}
	if k%2 == 0 {
		k = k + 1
	}
	for ; ; k += 2 {
		if isPrime(k) {
			return k
		}
	}
}

func isPrime(k int) bool {

	if k < 2 {
		return false
	}
	if k == 2 || k == 3 {
		return true
	}

	for i := 2; i < k/2+1; i++ {
		if k%i == 0 {
			return false
		}
	}
	return true
}
