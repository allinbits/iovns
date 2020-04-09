package keeper

import (
	"fmt"
	"github.com/iov-one/iovnsd"
	"github.com/iov-one/iovnsd/x/account/types"
	"github.com/iov-one/iovnsd/x/configuration"
	domain "github.com/iov-one/iovnsd/x/domain/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// expected keepers
type domainKeeper interface {
	GetDomain(ctx sdk.Context, domainName string) (domain domain.Domain, exists bool)
}

type configurationKeeper interface {
	GetConfig(ctx sdk.Context) configuration.Config
}

// Keeper of the account store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramspace types.ParamSubspace
	// external keepers
	ConfigKeeper configurationKeeper
	DomainKeeper domainKeeper
}

// NewKeeper creates a account keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramspace types.ParamSubspace) Keeper {
	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		paramspace: nil,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetAccount returns the account based on its name
// if the account does not exist it returns a zero value account type and false
func (k Keeper) GetAccount(ctx sdk.Context, accountName string) (types.Account, bool) {
	store := ctx.KVStore(k.storeKey)
	var item types.Account
	byteKey := []byte(accountName)
	accBytes := store.Get(byteKey)
	if len(accBytes) == 0 {
		return types.Account{}, false
	}
	k.cdc.MustUnmarshalBinaryBare(accBytes, &item)
	return item, true
}

// SetAccount sets the account
func (k Keeper) SetAccount(ctx sdk.Context, account types.Account) {
	accountKey := iovnsd.GetAccountKey(account.Domain, account.Name)
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(account)
	store.Set([]byte(accountKey), bz)
}

func (k Keeper) delete(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(key))
}
