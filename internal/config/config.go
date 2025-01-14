package config

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	RPCChainName          string `env:"RPC_CHAIN_NAME,default=ethereum"`
	RPCURL                string `env:"RPC_URL,default=http://localhost:8545"`
	RPCWSURL              string `env:"RPC_WS_URL,default=ws://localhost:8545"`
	EntryPointAddress     string `env:"ERC4337_ENTRYPOINT,required"`
	AccountFactoryAddress string `env:"ERC4337_ACCOUNT_FACTORY,required"`
	BundlerOriginHeader   string `env:"ERC4337_ORIGIN_HEADER,required"`
	ProfileAddress        string `env:"PROFILE_ADDRESS,required"`
	APIKEY                string `env:"API_KEY,required"`
	SentryURL             string `env:"SENTRY_URL"`
	PinataBaseURL         string `env:"PINATA_BASE_URL"`
	PinataAPIKey          string `env:"PINATA_API_KEY"`
	PinataAPISecret       string `env:"PINATA_API_SECRET"`
	DiscordURL            string `env:"DISCORD_URL,required"`
	PaymasterKey          string `env:"PAYMASTER_KEY,required"`
	DBSecret              string `env:"DB_SECRET,required"`
}

func New(ctx context.Context, envpath string) (*Config, error) {
	if envpath != "" {
		log.Default().Println("loading env from file: ", envpath)
		err := godotenv.Load(envpath)
		if err != nil {
			return nil, err
		}
	}

	cfg := &Config{}
	err := envconfig.Process(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
