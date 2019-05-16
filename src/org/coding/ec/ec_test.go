package ec

import (
	"testing"
	"math/rand"
	"time"
	"bytes"
	"fmt"
)

const (
	K = 5
	M = 4
	SZ = 1024 * 1024 * 32
)

const (
	KB  = 1024
	MB  = 1024 * KB
	GB  = 1024 * MB
)

func PrintRate(msg string, sz int, start time.Time)  {
	sec := time.Now().Sub(start).Seconds()
	rate := float64(sz) / sec

	if rate >= GB {
		fmt.Printf("%s: %f GB/s in %f seconds\n", msg, rate/GB, sec)
	} else if rate >= MB {
		fmt.Printf("%s: %f MB/s in %f seconds\n", msg, rate/MB, sec)
	} else {
		fmt.Printf("%s: %f KB/s in %f seconds\n", msg, rate/KB, sec)
	}
}

func Test_EC(t *testing.T)  {
	data := make([][]byte, K + M)
	in, out := [K][SZ]byte{}, [M][SZ]byte{}

	rand.Seed(time.Now().Unix())
	for i := 0; i < K; i += 1 {
		rand.Read(in[i][:])
		data[i] = in[i][:]
	}

	for i := 0; i < M; i += 1 {
		data[K + i] = out[i][:]
	}

	start := time.Now()
	Encode(data[:K], data[K:K+M])
	PrintRate("EC encoding", SZ * K, start)

	data1, row := [K][SZ]byte{}, [K]int{}
	recovery := make([][]byte, K)

	for i := 0; i < K; i += 1 {
		copy(data1[i][:], data[i + 4])
		recovery[i] = data1[i][:]
		row[i] = i + 4
	}

	start = time.Now()
	Decode(recovery, row[:])
	PrintRate("EC decoding", SZ * K, start)

	for i := 0; i < K; i += 1 {
		assert(bytes.Equal(data[i][:], recovery[i]))
	}
}
