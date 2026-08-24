package config

type Config struct {
	Rules RulesConfig `mapstructure:"rules"`
}

type RulesConfig struct {
	DebugLogging      DebugLoggingRuleConfig      `mapstructure:"debug_logging"`
	PlaintextPassword PlaintextPasswordRuleConfig `mapstructure:"plaintext_password"`
	PublicBind        PublicBindRuleConfig        `mapstructure:"public_bind"`
	DisabledTLS       DisabledTLSRuleConfig       `mapstructure:"disabled_tls"`
	WeakAlgorithm     WeakAlgorithmRuleConfig     `mapstructure:"weak_algorithm"`
}

type CommonRuleConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	Severity       string `mapstructure:"severity"`
	Recommendation string `mapstructure:"recommendation"`
}

type DebugLoggingRuleConfig struct {
	CommonRuleConfig `mapstructure:",squash"`
	Keys             []string `mapstructure:"keys"`
	UnsafeValues     []string `mapstructure:"unsafe_values"`
}

type PlaintextPasswordRuleConfig struct {
	CommonRuleConfig `mapstructure:",squash"`
	Keys             []string `mapstructure:"keys"`
}

type PublicBindRuleConfig struct {
	CommonRuleConfig `mapstructure:",squash"`
	Keys             []string `mapstructure:"keys"`
	UnsafeAddresses  []string `mapstructure:"unsafe_addresses"`
}

type DisabledTLSRuleConfig struct {
	CommonRuleConfig `mapstructure:",squash"`

	// Эти параметры опасны, когда имеют значение false.
	FalseKeys []string `mapstructure:"false_keys"`

	// Эти параметры опасны, когда имеют значение true.
	TrueKeys []string `mapstructure:"true_keys"`
}

type WeakAlgorithmRuleConfig struct {
	CommonRuleConfig `mapstructure:",squash"`
	Keys             []string `mapstructure:"keys"`
	UnsafeAlgorithms []string `mapstructure:"unsafe_algorithms"`
}
