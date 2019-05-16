package ec

func assert(cond bool)  {
	if !cond {
		panic("Assert fail")
	}
}

func cauchy(i, j, k int) uint8 {
	if i < k {
		if i == j {
			return GF_ONE
		}
		return GF_ZERO
	} else if i == k {
		return GF_ONE
	}
	// 1 / (i + j)
	return INVT[i ^ j]
}

// y = a * x
func vax(a uint8, x, y []byte)  {
	if a == GF_ONE {
		copy(y, x)
		return
	}
	for i := 0; i < len(x); i += 1 {
		y[i] = MT[a][x[i]]
	}
}

// y += a * x
func vaxpy(a uint8, x, y []byte)  {
	if a == GF_ZERO {
		return
	}
	for i := 0; i < len(x); i += 1 {
		y[i] ^= MT[a][x[i]]
	}
}

// y -= a * x
func vaxsy(a uint8, x, y []byte)  {
	vaxpy(a, x, y)
}

func Encode(in [][]byte, out [][]byte)  {
	k, m := len(in), len(out)
	assert(k <= MAX_K && m <= MAX_M)

	for i := 0; i < m; i += 1 {
		vax(cauchy(k + i, 0, k), in[0], out[i]);

		for j := 1; j < k; j += 1 {
			vaxpy(cauchy(k + i, j, k), in[j], out[i])
		}
	}
}

func swap(v [][]byte, i, j int)  {
	t := v[i]
	v[i] = v[j]
	v[j] = t
}

func pivot(A [][]uint8, i, k int) int {
	pivot := i

	for pivot < k && A[pivot][i] == GF_ZERO {
		pivot += 1
	}
	// not a reversible A
	assert(pivot < k)
	return pivot
}

func div(a, b uint8) uint8 {
	// a / b == a * (1/b)
	return MT[a][INVT[b]]
}

func solve(A [][]uint8, in [][]byte, k int) {
	// forward substitution
	for i := 0; i < k - 1; i += 1 {
		p := pivot(A, i, k)
		swap(A, i, p)
		swap(in, i, p)

		for j := i + 1; j < k; j += 1 {
			s := div(A[j][i], A[i][i])
			// A[j] -= s * A[i]
			vaxsy(s, A[i], A[j])
			// x[j] -= s * x[i]
			vaxsy(s, in[i], in[j])
		}
	}

	// backward substitution
	for i := k - 1; i > -1; i -= 1 {
		// x[i] /= A[i][i]
		vax(INVT[A[i][i]], in[i], in[i])
		A[i][i] = GF_ONE

		for j := i - 1; j > -1; j -= 1 {
			// x[j] -= A[j][i] * x[i]
			vaxsy(A[j][i], in[i], in[j])
			A[j][i] = GF_ZERO
		}
	}
}

func Decode(in [][]byte, row []int)  {
	k := len(in)
	assert(0 < k && k <= MAX_K && k == len(row))

	A := make([][]uint8, k)
	for i := 0; i < k; i += 1 {
		A[i] = make([]uint8, k)

		for j := 0; j < k; j += 1 {
			A[i][j] = cauchy(row[i], j, k)
		}
	}

	solve(A[:], in, k)
}