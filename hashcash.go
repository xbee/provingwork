package provingwork

import (
	"bytes"
	"fmt"
	"time"

	"math/big"

	"crypto/rand"
	"crypto/sha1"

	"encoding/base64"
	"encoding/binary"
)

// HashCash format:
// 1:20:20160927155710:somedatatovalidate::aW5ZdXJQcm90b2NvbHMh:VvJC
// version, zero bits, date, resource, extension (ignored), rand, counter

type HashCash struct {
	Counter  int64  `json:"counter"`
	Resource []byte `json:"resource"`

	*WorkOptions
}

func NewHashCash(resource []byte, opts ...*WorkOptions) *HashCash {
	hc := HashCash{
		Counter:  0,
		Resource: resource,
	}

	if (len(opts) != 0) {
		hc.WorkOptions = opts[0]
	} else {
		hc.WorkOptions = &WorkOptions{}
	}

	if hc.Timestamp == nil {
		t := time.Now()
		hc.Timestamp = &t
	}

	if hc.BitStrength == 0 {
		hc.BitStrength = DefaultBitStrength
	}

	if len(hc.Salt) == 0 {
		hc.Salt = make([]byte, DefaultSaltSize)
		rand.Read(hc.Salt)
	}

	return &hc
}

func (hc *HashCash) Check() bool {
	if (hc.ZeroCount() >= hc.BitStrength) {
		return true
	}
	return false
}

func (hc *HashCash) CounterBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, &hc.Counter)
	return buf.Bytes()
}

func (hc *HashCash) FindProof() {
	for {
		if hc.Check() {
			return
		}
		hc.Counter++
	}
}

func (hc *HashCash) String() string {
	return fmt.Sprintf(
		"1:%v:%v:%v:%v:%v:%v",
		hc.BitStrength,
		hc.Timestamp.Format("20060102150405"),
		string(hc.Resource),
		string(hc.Extension),
		base64.StdEncoding.EncodeToString(hc.Salt),
		base64.StdEncoding.EncodeToString(hc.CounterBytes()),
	)
}

func (hc *HashCash) ZeroCount() int {
	digest := sha1.Sum([]byte(hc.String()))
	digestHex := new(big.Int).SetBytes(digest[:])
	return ((sha1.Size * 8) - digestHex.BitLen())
}
