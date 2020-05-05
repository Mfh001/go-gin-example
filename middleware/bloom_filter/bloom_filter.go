package bloom_filter

import (
	"crypto/sha256"
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
)
import "github.com/spaolacci/murmur3"

type BloomFilter struct {
	FuncCount uint
	Len       uint
	Key       string
}

var Filter *BloomFilter

func init() {
	Filter = NewBloomFilter(4, 1000000)
}

func HashData(data []byte, seed uint) uint {
	shaData := sha256.Sum256(data)
	data = shaData[:]
	m := murmur3.New64WithSeed(uint32(seed))
	_, _ = m.Write(data)
	return uint(m.Sum64())
}

func (filter *BloomFilter) generateKey() {
	key := fmt.Sprintf("bloomFilter:n%d-l%d", filter.FuncCount, filter.Len)
	filter.Key = key
}

func (filter *BloomFilter) Put(data string) {
	for i := uint(0); i < filter.FuncCount; i++ {
		offset := HashData([]byte(data), i) % filter.Len
		_ = gredis.SetBit(filter.Key, offset, 1)
	}
}

func (filter *BloomFilter) Has(data string) bool {
	for i := uint(0); i < filter.FuncCount; i++ {
		offset := HashData([]byte(data), i) % filter.Len
		val, _ := gredis.GetBit(filter.Key, offset)
		if val != 1 {
			return false
		}
	}
	return true
}

func NewBloomFilter(funcCount uint, len uint) *BloomFilter {
	filter := &BloomFilter{
		FuncCount: funcCount,
		Len:       len,
	}
	filter.generateKey()
	//_ = gredis.SetBit(filter.Key, filter.Len, 0)
	return filter
}
