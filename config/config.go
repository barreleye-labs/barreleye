package config

var (
	BlockReward  = uint64(10)
	FaucetAmount = uint64(10)
)

type Config struct {
	PrivateKey string
}

func GetConfig(nodeName string) *Config {
	genesisConfig := CreateConfig("a2288db63c7016b815c55c1084c2491b8599834500408ba863ec379895373ae9")
	nayoungConfig := CreateConfig("c4e0f3f39c5438d2f7ba8b830f5a5538c6a63c752cb36fb1b91911539af01421")
	youngminConfig := CreateConfig("f2e1e4331b10c2b84a8ed58226398f5d11ee78052afa641d16851bd66bbdadb7")

	switch nodeName {
	case "genesis-node":
		return genesisConfig
	case "nayoung":
		return nayoungConfig
	case "youngmin":
		return youngminConfig
	}
	return nil
}

func CreateConfig(privateKey string) *Config {
	return &Config{
		PrivateKey: privateKey,
	}
}
