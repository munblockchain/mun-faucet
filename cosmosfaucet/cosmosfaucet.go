// Package cosmosfaucet is a faucet to request tokens for sdk accounts.
package cosmosfaucet

import (
	"context"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"faucet/cmd/config"

	chaincmdrunner "github.com/ignite/cli/ignite/pkg/chaincmd/runner"
	"gopkg.in/ezzarghili/recaptcha-go.v4"
)

const (
	// DefaultAccountName is the default account to transfer tokens from.
	DefaultAccountName = "faucet"

	// DefaultDenom is the default denomination to distribute.
	DefaultDenom = "uatom"

	// DefaultAmount specifies the default amount to transfer to an account
	// on each request.
	DefaultAmount = 10000000

	// DefaultMaxAmount specifies the maximum amount that can be tranffered to an
	// account in all times.
	DefaultMaxAmount = 100000000

	// DefaultLimitRefreshWindow specifies the time after which the max amount limit
	// is refreshed for an account [1 year]
	DefaultRefreshWindow = time.Hour * 24 * 36500
)

// Faucet represents a faucet.
type Faucet struct {
	// runner used to intereact with blockchain's binary to transfer tokens.
	runner chaincmdrunner.Runner

	// chainID is the chain id of the chain that faucet is operating for.
	chainID string

	// accountName to transfer tokens from.
	accountName string

	// accountMnemonic is the mnemonic of the account.
	accountMnemonic string

	// coinType registered coin type number for HD derivation (BIP-0044).
	coinType string

	// coins keeps a list of coins that can be distributed by the faucet.
	coins sdk.Coins

	// coinsMax is a denom-max pair.
	// it holds the maximum amounts of coins that can be sent to a single account.
	coinsMax map[string]uint64

	limitRefreshWindow time.Duration

	// openAPIData holds template data customizations for serving OpenAPI page & spec.
	openAPIData openAPIData

	configuration *config.Config
	captcha       recaptcha.ReCAPTCHA
}

// Option configures the faucetOptions.
type Option func(*Faucet)

// Account provides the account information to transfer tokens from.
// when mnemonic isn't provided, account assumed to be exists in the keyring.
func Account(name, mnemonic string, coinType string) Option {
	return func(f *Faucet) {
		f.accountName = name
		f.accountMnemonic = mnemonic
		f.coinType = coinType
	}
}

// Coin adds a new coin to coins list to distribute by the faucet.
// the first coin added to the list considered as the default coin during transfer requests.
//
// amount is the amount of the coin can be distributed per request.
// maxAmount is the maximum amount of the coin that can be sent to a single account.
// denom is denomination of the coin to be distributed by the faucet.
func Coin(amount, maxAmount uint64, denom string) Option {
	return func(f *Faucet) {
		f.coins = append(f.coins, sdk.NewCoin(denom, sdkmath.NewIntFromUint64(amount)))
		f.coinsMax[denom] = maxAmount
	}
}

// RefreshWindow adds the duration to refresh the transfer limit to the faucet
func RefreshWindow(refreshWindow time.Duration) Option {
	return func(f *Faucet) {
		f.limitRefreshWindow = refreshWindow
	}
}

// ChainID adds chain id to faucet. faucet will automatically fetch when it isn't provided.
func ChainID(id string) Option {
	return func(f *Faucet) {
		f.chainID = id
	}
}

// OpenAPI configures how to serve Open API page and and spec.
func OpenAPI(apiAddress string) Option {
	return func(f *Faucet) {
		f.openAPIData.APIAddress = apiAddress
	}
}

// New creates a new faucet with ccr (to access and use blockchain's CLI) and given options.
func New(ctx context.Context, ccr chaincmdrunner.Runner, conf *config.Config, options ...Option) (Faucet, error) {
	c, _ := recaptcha.NewReCAPTCHA(conf.ReCAPTCHA_ServerKey, recaptcha.V2, 10*time.Second) // for v2 API get your secret from https://www.google.com/recaptcha/admin

	f := Faucet{
		runner:        ccr,
		accountName:   DefaultAccountName,
		coinsMax:      make(map[string]uint64),
		openAPIData:   openAPIData{"Blockchain", "http://localhost:1317"},
		configuration: conf,
		captcha:       c,
	}

	for _, apply := range options {
		apply(&f)
	}

	if len(f.coins) == 0 {
		Coin(DefaultAmount, DefaultMaxAmount, DefaultDenom)(&f)
	}

	if f.limitRefreshWindow == 0 {
		RefreshWindow(DefaultRefreshWindow)(&f)
	}

	// import the account if mnemonic is provided.
	if f.accountMnemonic != "" {
		_, err := f.runner.AddAccount(ctx, f.accountName, f.accountMnemonic, f.coinType)
		if err != nil && err != chaincmdrunner.ErrAccountAlreadyExists {
			return Faucet{}, err
		}
	}

	if f.chainID == "" {
		status, err := f.runner.Status(ctx)
		if err != nil {
			return Faucet{}, err
		}

		f.chainID = status.ChainID
		f.openAPIData.ChainID = status.ChainID
	}

	return f, nil
}
