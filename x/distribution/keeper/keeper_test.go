package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

func TestSetWithdrawAddr(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInputDefault(t, false, 1000)

	keeper.SetWithdrawAddrEnabled(ctx, false)

	err := keeper.SetWithdrawAddr(ctx, delAddr1, delAddr2)
	require.NotNil(t, err)

	keeper.SetWithdrawAddrEnabled(ctx, true)

	err = keeper.SetWithdrawAddr(ctx, delAddr1, delAddr2)
	require.Nil(t, err)
}

func TestWithdrawValidatorCommission(t *testing.T) {
	ctx, ak, keeper, _, _ := CreateTestInputDefault(t, false, 1000)

	// set zero outstanding rewards
	keeper.SetOutstandingRewards(ctx, types.OutstandingRewards{})

	// check initial balance
	balance := ak.GetAccount(ctx, sdk.AccAddress(valOpAddr3)).GetCoins()
	expTokens := sdk.TokensFromTendermintPower(1000)
	require.Equal(t, sdk.Coins{
		{"stake", sdk.TokensFromTendermintPower(1000)},
	}, balance)

	// set commission
	keeper.SetValidatorAccumulatedCommission(ctx, valOpAddr3, sdk.DecCoins{
		{"mytoken", sdk.NewDec(5).Quo(sdk.NewDec(4))},
		{"stake", sdk.NewDec(3).Quo(sdk.NewDec(2))},
	})

	// withdraw commission
	keeper.WithdrawValidatorCommission(ctx, valOpAddr3)

	// check balance increase
	balance = ak.GetAccount(ctx, sdk.AccAddress(valOpAddr3)).GetCoins()
	require.Equal(t, sdk.Coins{
		{"mytoken", sdk.NewInt(1)},
		{"stake", expTokens.AddRaw(1)},
	}, balance)

	// check remainder
	remainder := keeper.GetValidatorAccumulatedCommission(ctx, valOpAddr3)
	require.Equal(t, sdk.DecCoins{
		{"mytoken", sdk.NewDec(1).Quo(sdk.NewDec(4))},
		{"stake", sdk.NewDec(1).Quo(sdk.NewDec(2))},
	}, remainder)

	require.True(t, true)
}
