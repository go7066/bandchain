package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/chain/v2/testing/testapp"
)

func TestCheckSelfReporter(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// Owner must always be a reporter of himself.
	require.True(t, k.IsReporter(ctx, testapp.Owner.ValAddress, testapp.Owner.Address))
}

func TestAddReporter(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// Before we do anything, Bob and Carol must not be a reporter of Alice.
	require.False(t, k.IsReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address))
	// Adds Bob as a reporter Alice. IsReporter should return true for Bob, false for Carol.
	err := k.AddReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address)
	require.Nil(t, err)
	require.True(t, k.IsReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address))
	require.False(t, k.IsReporter(ctx, testapp.Alice.ValAddress, testapp.Carol.Address))
	// We should get an error if we try to add Bob again.
	err = k.AddReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address)
	require.NotNil(t, err)
}

func TestRemoveReporter(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// Removing Bob from Alice's reporter list should error as Bob is not a reporter of Alice.
	err := k.RemoveReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address)
	require.NotNil(t, err)
	// Adds Bob as the reporter of Alice. We now should be able to remove Bob, but not Carol.
	err = k.AddReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address)
	require.Nil(t, err)
	err = k.RemoveReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address)
	require.Nil(t, err)
	err = k.RemoveReporter(ctx, testapp.Alice.ValAddress, testapp.Carol.Address)
	require.NotNil(t, err)
	// By the end of everything, no one should be a reporter of Alice.
	require.False(t, k.IsReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address))
	require.False(t, k.IsReporter(ctx, testapp.Alice.ValAddress, testapp.Carol.Address))
}

func TestGetReporters(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// Initially, only Alice should be the reporter of Alice.
	reporters := k.GetReporters(ctx, testapp.Alice.ValAddress)
	require.Equal(t, 1, len(reporters))
	require.Contains(t, reporters, testapp.Alice.Address)
	// After we add Bob and Carol, they should also appear in GetReporters.
	err := k.AddReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address)
	require.NoError(t, err)
	err = k.AddReporter(ctx, testapp.Alice.ValAddress, testapp.Carol.Address)
	require.NoError(t, err)
	reporters = k.GetReporters(ctx, testapp.Alice.ValAddress)
	require.Equal(t, 3, len(reporters))
	require.Contains(t, reporters, testapp.Alice.Address)
	require.Contains(t, reporters, testapp.Bob.Address)
	require.Contains(t, reporters, testapp.Carol.Address)
}

func TestGetAllReporters(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// Initially, only validators should be reporters of themselves
	reporters := k.GetAllReporters(ctx)
	expectedReporterMap := map[string][]string{
		sdk.ValAddress(testapp.Validators[0].Address).String(): {testapp.Validators[0].Address.String()},
		sdk.ValAddress(testapp.Validators[1].Address).String(): {testapp.Validators[1].Address.String()},
		sdk.ValAddress(testapp.Validators[2].Address).String(): {testapp.Validators[2].Address.String()},
	}
	require.Equal(t, len(expectedReporterMap), len(reporters))
	for _, reportersPerValidator := range reporters {
		require.Contains(t, expectedReporterMap, reportersPerValidator.Validator)
		for _, reporter := range reportersPerValidator.Reporters {
			require.Contains(t, expectedReporterMap[reportersPerValidator.Validator], reporter)
		}
	}

	// After Alice, Bob, and Carol are added, they should be included in result of GetAllReporters
	err := k.AddReporter(ctx, testapp.Validators[0].ValAddress, testapp.Alice.Address)
	require.NoError(t, err)
	err = k.AddReporter(ctx, testapp.Validators[0].ValAddress, testapp.Bob.Address)
	require.NoError(t, err)
	err = k.AddReporter(ctx, testapp.Validators[2].ValAddress, testapp.Carol.Address)
	require.NoError(t, err)

	reporters = k.GetAllReporters(ctx)
	expectedReporterMap = map[string][]string{
		sdk.ValAddress(testapp.Validators[0].Address).String(): {testapp.Validators[0].Address.String(), testapp.Bob.Address.String(), testapp.Alice.Address.String()},
		sdk.ValAddress(testapp.Validators[1].Address).String(): {testapp.Validators[1].Address.String()},
		sdk.ValAddress(testapp.Validators[2].Address).String(): {testapp.Validators[2].Address.String(), testapp.Carol.Address.String()},
	}

	require.Equal(t, len(expectedReporterMap), len(reporters))
	for _, reportersPerValidator := range reporters {
		require.Contains(t, expectedReporterMap, reportersPerValidator.Validator)
		for _, reporter := range reportersPerValidator.Reporters {
			require.Contains(t, expectedReporterMap[reportersPerValidator.Validator], reporter)
		}
	}
}
