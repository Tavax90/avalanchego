// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package sampler

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"golang.org/x/exp/slices"

	safemath "github.com/ava-labs/avalanchego/utils/math"
)

var (
	weightedWithoutReplacementSamplers = []struct {
		name    string
		sampler WeightedWithoutReplacement
	}{
		{
			name: "generic with replacer and best",
			sampler: &weightedWithoutReplacementGeneric{
				u: &uniformReplacer{
					rng: globalRNG,
				},
				w: &weightedBest{
					samplers: []Weighted{
						&weightedArray{},
						&weightedHeap{},
						&weightedUniform{
							maxWeight: 1024,
						},
					},
					benchmarkIterations: 30,
				},
			},
		},
	}
	weightedWithoutReplacementTests = []struct {
		name string
		test func(*testing.T, WeightedWithoutReplacement)
	}{
		{
			name: "initialize overflow",
			test: WeightedWithoutReplacementInitializeOverflowTest,
		},
		{
			name: "out of range",
			test: WeightedWithoutReplacementOutOfRangeTest,
		},
		{
			name: "empty without weight",
			test: WeightedWithoutReplacementEmptyWithoutWeightTest,
		},
		{
			name: "empty",
			test: WeightedWithoutReplacementEmptyTest,
		},
		{
			name: "singleton",
			test: WeightedWithoutReplacementSingletonTest,
		},
		{
			name: "with zero",
			test: WeightedWithoutReplacementWithZeroTest,
		},
		{
			name: "distribution",
			test: WeightedWithoutReplacementDistributionTest,
		},
	}
)

func TestAllWeightedWithoutReplacement(t *testing.T) {
	for _, s := range weightedWithoutReplacementSamplers {
		for _, test := range weightedWithoutReplacementTests {
			t.Run(fmt.Sprintf("sampler %s test %s", s.name, test.name), func(t *testing.T) {
				test.test(t, s.sampler)
			})
		}
	}
}

func WeightedWithoutReplacementInitializeOverflowTest(
	t *testing.T,
	s WeightedWithoutReplacement,
) {
	err := s.Initialize([]uint64{1, math.MaxUint64})
	require.ErrorIs(t, err, safemath.ErrOverflow)
}

func WeightedWithoutReplacementOutOfRangeTest(
	t *testing.T,
	s WeightedWithoutReplacement,
) {
	require := require.New(t)

	require.NoError(s.Initialize([]uint64{1}))

	_, err := s.Sample(2)
	require.ErrorIs(err, ErrOutOfRange)
}

func WeightedWithoutReplacementEmptyWithoutWeightTest(
	t *testing.T,
	s WeightedWithoutReplacement,
) {
	require := require.New(t)

	require.NoError(s.Initialize(nil))

	indices, err := s.Sample(0)
	require.NoError(err)
	require.Empty(indices)
}

func WeightedWithoutReplacementEmptyTest(
	t *testing.T,
	s WeightedWithoutReplacement,
) {
	require := require.New(t)

	require.NoError(s.Initialize([]uint64{1}))

	indices, err := s.Sample(0)
	require.NoError(err)
	require.Empty(indices)
}

func WeightedWithoutReplacementSingletonTest(
	t *testing.T,
	s WeightedWithoutReplacement,
) {
	require := require.New(t)

	require.NoError(s.Initialize([]uint64{1}))

	indices, err := s.Sample(1)
	require.NoError(err)
	require.Equal([]int{0}, indices)
}

func WeightedWithoutReplacementWithZeroTest(
	t *testing.T,
	s WeightedWithoutReplacement,
) {
	require := require.New(t)

	require.NoError(s.Initialize([]uint64{0, 1}))

	indices, err := s.Sample(1)
	require.NoError(err)
	require.Equal([]int{1}, indices)
}

func WeightedWithoutReplacementDistributionTest(
	t *testing.T,
	s WeightedWithoutReplacement,
) {
	require := require.New(t)

	require.NoError(s.Initialize([]uint64{1, 1, 2}))

	indices, err := s.Sample(4)
	require.NoError(err)

	slices.Sort(indices)
	require.Equal([]int{0, 1, 2, 2}, indices)
}
