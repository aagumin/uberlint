package timefield

type BadConfig struct {
	Timeout int `json:"timeout"` // want "include time unit in serialized numeric field name"
	TTL     int `yaml:"ttl"`     // want "include time unit in serialized numeric field name"
}

type GoodConfig struct {
	TimeoutMillis int `json:"timeoutMillis"`
	TTLSeconds    int `yaml:"ttlSeconds"`
	RetryCount    int `json:"retryCount"`
}
