# go-leb128
> Little Endian Base 128, a variable-length code compression for arbitrarily large integer in a small number of bytes.

# Usage 
The bytes payload is typed in this library. So instead of using `[]byte` we are using

* `SLeb128` for Signed LEB128 encoding
* `ULeb128` for Unsigned LEB128 encoding

So for unsigned LEB128 encoding it should be written in

```go
sleb624485 := ULeb128{[]byte{0xE5, 0x8E, 0x26}}
```

And it cannot be consumed by signed LEB128 decoder (typed)

```go
DecodeSLeb128(sleb624485) // this will not compile
DecodeULeb128(sleb624485) // compiles !
```

### Encoding

**EncodeFromUint64(uint64) ULeb128**

Encode `uint64` typed integer into `ULeb128`

**EncodeFromInt64(int64) SLeb128**

Encode `int64` typed integer into `SLeb128`


### Decoding

**DecodeULeb128(ULeb128) (uint64, error)**

Decode `ULeb128` to uint64 with error flag

**DecodeSLeb128(SLeb128) (int64, error)**

Decode `SLeb128` to int64 with error flag

> If the error is not equal to nil (means error occurred). The result is always 0

### Utilities

**AppendSLeb128(x, y SLeb128) SLeb128**

Append two `Sleb128` into one `SLeb128`

**AppendULeb128(x, y ULeb128) ULeb128**

Append two `Uleb128` into one `ULeb128`
