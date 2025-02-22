package bytecode

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/SpecDrivenDesign/lql/pkg/tokens"
)

// ByteCodeReader reads tokens from a binary-encoded byte slice.
type ByteCodeReader struct {
	data []byte
	pos  int
}

// NewByteCodeReader creates a new ByteCodeReader.
func NewByteCodeReader(data []byte) *ByteCodeReader {
	return &ByteCodeReader{
		data: data,
		pos:  0,
	}
}

// NextToken decodes the next token.
func (b *ByteCodeReader) NextToken() (tokens.Token, error) {
	if b.pos >= len(b.data) {
		return tokens.Token{Type: tokens.TokenEof, Literal: ""}, nil
	}

	// Read token type byte.
	tokenTypeByte := b.data[b.pos]
	b.pos++
	tokenType, ok := ByteToTokenType[tokenTypeByte]
	if !ok {
		return tokens.Token{Type: tokens.TokenIllegal, Literal: ""}, fmt.Errorf("unknown token type code: %v", tokenTypeByte)
	}

	var literal string
	// If the token has a fixed literal, use that.
	if fixed, isFixed := tokens.FixedTokenLiterals[tokenType]; isFixed {
		literal = fixed
	} else {
		// Otherwise, read a length-prefixed literal.
		if b.pos+1 > len(b.data) {
			return tokens.Token{Type: tokens.TokenIllegal, Literal: ""}, fmt.Errorf("unexpected end of data reading literal length")
		}
		length := b.data[b.pos]
		b.pos++
		if b.pos+int(length) > len(b.data) {
			return tokens.Token{Type: tokens.TokenIllegal, Literal: ""}, fmt.Errorf("unexpected end of data reading literal")
		}
		literal = string(b.data[b.pos : b.pos+int(length)])
		b.pos += int(length)
	}

	// Construct the token. Note: line/column info isn't preserved here.
	return tokens.Token{
		Type:    tokenType,
		Literal: literal,
		Line:    -1,
		Column:  -1,
	}, nil
}

// NewByteCodeReaderFromSignedData verifies the RSA signature over the token data
// and returns a ByteCodeReader if the signature is valid.
func NewByteCodeReaderFromSignedData(data []byte, pub *rsa.PublicKey) (*ByteCodeReader, error) {
	sigSize := pub.Size() // RSA signature size in bytes.
	if len(data) < len(tokens.HeaderMagic)+4+sigSize {
		return nil, fmt.Errorf("data too short to contain valid signed tokens")
	}

	if string(data[:len(tokens.HeaderMagic)]) != tokens.HeaderMagic {
		return nil, fmt.Errorf("invalid header magic; expected %s", tokens.HeaderMagic)
	}
	pos := len(tokens.HeaderMagic)

	// Read the 4-byte little-endian length of tokenData.
	tokenDataLength := binary.LittleEndian.Uint32(data[pos : pos+4])
	pos += 4

	expectedLength := len(tokens.HeaderMagic) + 4 + int(tokenDataLength) + sigSize
	if len(data) != expectedLength {
		return nil, fmt.Errorf("data length mismatch: expected %d bytes, got %d", expectedLength, len(data))
	}

	tokenData := data[pos : pos+int(tokenDataLength)]
	pos += int(tokenDataLength)
	signature := data[pos : pos+sigSize]

	// Compute SHA256 hash over tokenData.
	hash := sha256.Sum256(tokenData)
	// Verify the RSA signature.
	if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], signature); err != nil {
		return nil, fmt.Errorf("invalid signature: %v", err)
	}

	return NewByteCodeReader(tokenData), nil
}

// And a reverse mapping to convert a byte code back to a TokenType.
var ByteToTokenType = func() map[byte]tokens.TokenType {
	m := make(map[byte]tokens.TokenType)
	for tt, b := range tokens.TokenTypeToByte {
		m[b] = tt
	}
	return m
}()
