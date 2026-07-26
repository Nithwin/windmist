package rag

import (
	"encoding/binary"
	"math"
)

// Vector represents a float32 embedding vector.
type Vector = []float32

// CosineSimilarity computes the cosine similarity between two vectors.
// Returns a value between -1 and 1, where 1 means identical direction.
func CosineSimilarity(a, b Vector) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float32
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

// EncodeVector serializes a float32 vector to bytes for SQLite BLOB storage.
func EncodeVector(v Vector) []byte {
	buf := make([]byte, len(v)*4)
	for i, f := range v {
		binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(f))
	}
	return buf
}

// DecodeVector deserializes bytes from a SQLite BLOB back to a float32 vector.
func DecodeVector(buf []byte) Vector {
	if len(buf)%4 != 0 {
		return nil
	}
	v := make(Vector, len(buf)/4)
	for i := range v {
		v[i] = math.Float32frombits(binary.LittleEndian.Uint32(buf[i*4:]))
	}
	return v
}

// Normalize normalizes a vector to unit length (L2 normalization).
func Normalize(v Vector) Vector {
	var sum float32
	for _, f := range v {
		sum += f * f
	}
	if sum == 0 {
		return v
	}
	norm := float32(math.Sqrt(float64(sum)))
	out := make(Vector, len(v))
	for i, f := range v {
		out[i] = f / norm
	}
	return out
}
