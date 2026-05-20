package infra

import (
	"context"
	"hash/fnv"
	"math"

	"github.com/cloudwego/eino/components/embedding"
)

const fakeEmbedDim = 128

// FakeEmbedder produces deterministic unit vectors for local dev/tests without Ark.
type FakeEmbedder struct{}

func NewFakeEmbedder() embedding.Embedder {
	return &FakeEmbedder{}
}

func (f *FakeEmbedder) EmbedStrings(ctx context.Context, texts []string, opts ...embedding.Option) ([][]float64, error) {
	_ = ctx
	_ = opts
	out := make([][]float64, len(texts))
	for i, t := range texts {
		out[i] = hashToVector(t, fakeEmbedDim)
	}
	return out, nil
}

func hashToVector(text string, dim int) []float64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(text))
	seed := h.Sum64()
	vec := make([]float64, dim)
	var norm float64
	for i := range vec {
		seed = seed*6364136223846793005 + 1
		v := float64(int64(seed)%1000) / 500.0
		vec[i] = v
		norm += v * v
	}
	if norm > 0 {
		scale := 1.0 / math.Sqrt(norm)
		for i := range vec {
			vec[i] *= scale
		}
	}
	return vec
}

// FakeEmbedDim is the vector dimension used by FakeEmbedder and default Redis index.
func FakeEmbedDim() int { return fakeEmbedDim }

// FakeEmbedVector returns a deterministic unit vector for one query string.
func FakeEmbedVector(text string) []float32 {
	vec := hashToVector(text, fakeEmbedDim)
	out := make([]float32, len(vec))
	for i, v := range vec {
		out[i] = float32(v)
	}
	return out
}
