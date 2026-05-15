package timefield

type BadConfig struct {
	Timeout int     `json:"timeout"` // want "include time unit in serialized numeric field name"
	TTL     int     `yaml:"ttl"`     // want "include time unit in serialized numeric field name"
	Delay   float64 `json:"delay"`   // want "include time unit in serialized numeric field name"
	Expires int64   `json:"expires"` // want "include time unit in serialized numeric field name"
}

type GoodConfig struct {
	TimeoutMillis int     `json:"timeoutMillis"`
	TTLSeconds    int     `yaml:"ttlSeconds"`
	RetryCount    int     `json:"retryCount"`
	DelayMS       float64 `json:"delayMs"`
	IntervalHours int64   `yaml:"intervalHours"`
	Timeout       string  `json:"timeout"`
	Duration      int     `db:"duration"`
}
