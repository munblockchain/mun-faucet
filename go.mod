module faucet

go 1.16

require (
	cosmossdk.io/math v1.0.0-beta.3
	github.com/cosmos/cosmos-sdk v0.46.2
	github.com/gorilla/mux v1.8.0
	github.com/ignite/cli v0.24.1
	github.com/joho/godotenv v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.8.2
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.8.0
	github.com/tendermint/starport v0.19.5
	gopkg.in/ezzarghili/recaptcha-go.v4 v4.3.0 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
