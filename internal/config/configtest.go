package config

// type ConfigTest struct {
// 	DBName  string `mapstructure:"DATABASE_NAME"`
// 	DBPass  string `mapstructure:"DATABASE_PASSWORD"`
// 	DBUser  string `mapstructure:"DATABASE_USER"`
// 	DBHost  string `mapstructure:"DATABASE_HOST"`
// 	DBPort  string `mapstructure:"DATABASE_PORT"`
// 	AppPort string `mapstructure:"PORT"`
// }

func LoadTestConfig() *Config {
	dbUser := "testuser"
	dbPassword := "testpassword"
	dbHost := "localhost"
	dbPort := "5433" 
	dbName := "testdb"

	return &Config{
		DBName: dbName,
		DBPass: dbPassword,
		DBUser: dbUser,
		DBHost: dbHost,
		DBPort: dbPort,

	}
}
