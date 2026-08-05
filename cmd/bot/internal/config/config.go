package config

import (
	"time"
)

type Config struct {
	LogLevel         string                      `yaml:"log_level" required:"true"`
	Tracing          Tracing                     `yaml:"tracing" required:"true"`
	Bot              Bot                         `yaml:"bot" required:"true"`
	Postgres         Postgres                    `yaml:"postgres" required:"true"`
	BoxRewardPercent map[string]BoxRewardPercent `yaml:"box_reward_percent" required:"true"`
	BoxCosts         map[string]BoxCost          `yaml:"box_costs" required:"true"`
	BoxWaitPeriod    map[string]time.Duration    `yaml:"box_wait_period" required:"true"`
}

func (c *Config) Level() string {
	return c.LogLevel
}

func (c *Config) ServiceName() string {
	return c.Tracing.ServiceName
}

func (c *Config) TracingEndpoint() string {
	return c.Tracing.Endpoint
}

type Tracing struct {
	Endpoint    string `yaml:"endpoint" required:"true"`
	ServiceName string `yaml:"service_name" required:"true"`
}

type Bot struct {
	Token        string `yaml:"token" required:"true"`
	WebHookToken string `yaml:"webhook_token"`
}

type Postgres struct {
	Connection string `yaml:"connection" required:"true"`
}

type BoxRewardPercent struct {
	Legendary int `yaml:"legendary" required:"true"`
	Epic      int `yaml:"epic" required:"true"`
	Rare      int `yaml:"rare" required:"true"`
}

type BoxCost struct {
	Amount int `yaml:"amount" required:"true"`
}
