package config

import "fmt"

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
	)
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf(":%d", c.Server.Port)
}

func (c *Config) GetPostgresDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)
}

func (c *Config) GetRedisDSN() string {
	if c.Redis.Password != "" {
		return fmt.Sprintf("redis://:%s@%s:%d/%d", c.Redis.Password, c.Redis.Host, c.Redis.Port, c.Redis.DB)
	}
	return fmt.Sprintf("redis://%s:%d/%d", c.Redis.Host, c.Redis.Port, c.Redis.DB)
}
