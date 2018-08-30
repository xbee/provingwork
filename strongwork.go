package provingwork

import (
	"bytes"
	"fmt"

	"crypto/sha256"
	"math/big"

	"encoding/base64"
	"encoding/binary"
	"encoding/json"
)

type StrongWork struct {
	Counter  int64  `json:"counter"`
	Resource []byte `json:"resource"`

	*WorkOptions
}

// An alias type that won't have any of functions (mostly to avoid an infinite
// loop with the overidden MarshalJSON function)
type RawStrongWork StrongWork

// This is a special version of the StrongWork that has the types we want to
// be importing / exporting.
type StrongWorkExt struct {
	Timestamp int64 `json:"timestamp"`

	*RawStrongWork
}

func (wo StrongWork) MarshalJSON() ([]byte, error) {
	woe := StrongWorkExt{RawStrongWork: (*RawStrongWork)(&wo)}

	if wo.Timestamp != nil {
		woe.Timestamp = wo.Timestamp.Unix()
	}

	return json.Marshal(woe)
}

func (wo StrongWork) UnmarshalJSON(data []byte) error {
	woe := StrongWorkExt{RawStrongWork: (*RawStrongWork)(&wo)}

	if err := json.Unmarshal(data, woe); err != nil {
		return err
	}

	return nil
}

func NewStrongWork(resource []byte, opts ...*WorkOptions) *StrongWork {
	sw := StrongWork{Resource: resource}

	if len(opts) != 0 {
		sw.WorkOptions = opts[0]
	} else {
		sw.WorkOptions = &WorkOptions{}
	}

	setDefaultWorkOptions(sw.WorkOptions)

	return &sw
}

func (sw StrongWork) Check() bool {
	if sw.ZeroCount() >= sw.BitStrength {
		return true
	}
	return false
}

func (sw StrongWork) ContentBytes() []byte {
	var buf bytes.Buffer

	buf.Write(sw.Resource)
	buf.Write(sw.Salt)

	ts := sw.Timestamp.Unix()
	binary.Write(&buf, binary.BigEndian, ts)
	binary.Write(&buf, binary.BigEndian, sw.Counter)

	return buf.Bytes()
}

func (sw StrongWork) ContentHash() []byte {
	result := sha256.Sum256(sw.ContentBytes())
	return result[:]
}

func (sw StrongWork) CounterBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, sw.Counter)
	return buf.Bytes()
}

func (sw StrongWork) String() string {
	return fmt.Sprintf(
		"1:%v:%v:%v:%v:%v:%v",
		sw.BitStrength,
		sw.Timestamp.Format("20060102150405"),
		string(sw.Resource),
		string(sw.Extension),
		base64.StdEncoding.EncodeToString(sw.Salt),
		base64.StdEncoding.EncodeToString(sw.CounterBytes()),
	)
}

func (sw *StrongWork) FindProof() {
	for {
		if sw.Check() {
			return
		}
		sw.Counter++
	}
}

func (sw StrongWork) ZeroCount() int {
	digest := sha256.Sum256(sw.ContentBytes())
	digestHex := new(big.Int).SetBytes(digest[:])
	return ((sha256.Size * 8) - digestHex.BitLen())
}
