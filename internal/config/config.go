package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Пакет config нужен для чтения конфига самой утилиты echelon-utility
// Нужен для определения правил (например, какие алгоритмы считаем ненадежными и тп)

// Load считывает переменные из файла config.yaml, находящегося по пути cli/cmd/echelon-utility
func Load() (*Config, error) {
	v := viper.New()
	setDefaults(v)

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.AddConfigPath("./cli/cmd/echelon-utility")
	v.AddConfigPath(".")

	// Если конфиг не найден - пользуемся default значениями
	// Если найден, но возникла ошибка чтения - возвращаем ошибку
	if err := v.ReadInConfig(); err != nil {
		var notFoundErr viper.ConfigFileNotFoundError

		if !errors.As(err, &notFoundErr) {
			return nil, fmt.Errorf("read config error: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config error: %w", err)
	}

	return &cfg, nil
}

// setDefaults - выставляет дефолтные значения, если конфига нет
func setDefaults(v *viper.Viper) {

	// Логирование
	v.SetDefault("rules.debug_logging.enabled", true)
	v.SetDefault("rules.debug_logging.severity", "LOW")
	v.SetDefault("rules.debug_logging.recommendation",
		"Используйте менее подробный уровень логирования (info+).")
	v.SetDefault("rules.debug_logging.keys", []string{
		"level",
		"log_level",
	})
	v.SetDefault("rules.debug_logging.unsafe_values", []string{
		"debug",
		"trace",
	})

	// Пароли
	v.SetDefault("rules.plaintext_password.enabled", true)
	v.SetDefault("rules.plaintext_password.severity", "HIGH")
	v.SetDefault("rules.plaintext_password.recommendation",
		"Не храните пароль в открытом виде. Используйте защищённое хранилище секретов.")
	v.SetDefault("rules.plaintext_password.keys", []string{
		"password",
		"passwd",
		"pwd",
		"db_password",
	})

	// Public bind
	v.SetDefault("rules.public_bind.enabled", true)
	v.SetDefault("rules.public_bind.severity", "MEDIUM")
	v.SetDefault("rules.public_bind.recommendation",
		"Ограничьте адрес прослушивания или доступ с помощью сетевых правил.")
	v.SetDefault("rules.public_bind.keys", []string{
		"host",
		"bind",
		"listen_address",
	})
	v.SetDefault("rules.public_bind.unsafe_addresses", []string{
		"0.0.0.0",
		"::",
	})

	// TLS
	v.SetDefault("rules.disabled_tls.enabled", true)
	v.SetDefault("rules.disabled_tls.severity", "HIGH")
	v.SetDefault("rules.disabled_tls.recommendation",
		"Включите TLS и проверку сертификата.")
	v.SetDefault("rules.disabled_tls.false_keys", []string{
		"tls",
		"tls_enabled",
		"tls_verify",
	})
	v.SetDefault("rules.disabled_tls.true_keys", []string{
		"insecure",
		"skip_tls_verify",
		"insecure_skip_verify",
	})

	// Слабые алгоритмы
	v.SetDefault("rules.weak_algorithm.enabled", true)
	v.SetDefault("rules.weak_algorithm.severity", "HIGH")
	v.SetDefault("rules.weak_algorithm.recommendation",
		"Замените алгоритм на более безопасный.")
	v.SetDefault("rules.weak_algorithm.keys", []string{
		"algorithm",
		"hash_algorithm",
		"digest_algorithm",
		"digest-algorithm",
	})
	v.SetDefault("rules.weak_algorithm.unsafe_algorithms", []string{
		"md5",
		"sha1",
		"des",
	})

	// Права доступа к файлу
	v.SetDefault("rules.file_permissions.enabled", true)
	v.SetDefault("rules.file_permissions.severity", "MEDIUM")
	v.SetDefault("rules.file_permissions.recommendation",
		"Уберите права на запись для группы и остальных пользователей.")
}
