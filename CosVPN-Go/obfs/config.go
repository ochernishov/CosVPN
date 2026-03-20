package obfs

// ObfsConfig хранит параметры обфускации
type ObfsConfig struct {
	Key  [16]byte // 128-бит ключ для XOR заголовков
	Mode string   // "auto", "direct", "tls"
}

func DefaultConfig() ObfsConfig {
	return ObfsConfig{Mode: "direct"}
}

func (c *ObfsConfig) IsEnabled() bool {
	var zero [16]byte
	return c.Key != zero
}
