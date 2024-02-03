package config

func Default() *GlobalRateLimiter {
	return new(GlobalRateLimiter)
}

type GlobalRateLimiter struct {
	RateLimiter RateLimiter `yaml:"rateLimiter"`
}
type Limits struct {
	Count int `yaml:"count"`
	Rate  int `yaml:"rate"`
}
type RateLimiter struct {
	Limits Limits `yaml:"limits"`
}
