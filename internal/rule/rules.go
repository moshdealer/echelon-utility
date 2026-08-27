package rule

import (
	"echelon-utility/internal/config"
	"fmt"
	"strings"
)

/* Здесь описаны все необходимых 5 правил из ТЗ
Для добавления нового правила необходимо реализовать новую структуру и метод Check + фабрику для него в этом модуле
Также дополнить конфиг для правила
На текущий момент не было цели сделать движок, который позволял бы единообразно описывать n-количество правил
Сейчас реализация позволяет гибко настраивать проверяемые ключи и значения в рамках изначальных 5 правил
Если нужно будет добавить новое правило, то придется реализовать его ниже
В целом, требования по масштабированию соблюдены
*/

type baseRule struct {
	enabled        bool
	severity       string
	recommendation string
}
type LogRule struct {
	Base         baseRule
	Keys         map[string]struct{}
	UnsafeValues map[string]struct{}
}

func newBaseRule(cfg config.CommonRuleConfig) baseRule {
	return baseRule{
		enabled:        cfg.Enabled,
		severity:       cfg.Severity,
		recommendation: cfg.Recommendation,
	}
}

func NewLogRule(cfgRule config.DebugLoggingRuleConfig) *LogRule {
	keyMap := make(map[string]struct{})
	valueMap := make(map[string]struct{})

	for _, key := range cfgRule.Keys {
		keyMap[strings.ToLower(key)] = struct{}{}
	}

	for _, value := range cfgRule.UnsafeValues {
		valueMap[strings.ToLower(value)] = struct{}{}
	}

	return &LogRule{
		Base:         newBaseRule(cfgRule.CommonRuleConfig),
		Keys:         keyMap,
		UnsafeValues: valueMap,
	}
}

func (r *LogRule) Check(path string, key string, value any) *Finding {
	ruleFlag := false
	stringValue, ok := value.(string)
	if !ok {
		return nil
	}

	if _, ok := r.Keys[strings.ToLower(key)]; ok {
		if _, ok := r.UnsafeValues[strings.ToLower(stringValue)]; ok {
			ruleFlag = true
		}
	}

	if ruleFlag {
		return &Finding{
			Severity:       r.Base.severity,
			Rule:           "Правило логирования",
			Path:           path,
			Message:        fmt.Sprintf("Логирование в режиме %q", stringValue),
			Recommendation: r.Base.recommendation,
		}
	}
	return nil
}

type PassRule struct {
	Base baseRule
	Keys map[string]struct{}
}

func NewPassRule(cfgRule config.PlaintextPasswordRuleConfig) *PassRule {
	keyMap := make(map[string]struct{})

	for _, key := range cfgRule.Keys {
		keyMap[strings.ToLower(key)] = struct{}{}
	}

	return &PassRule{
		Base: newBaseRule(cfgRule.CommonRuleConfig),
		Keys: keyMap,
	}
}

func (r *PassRule) Check(path string, key string, value any) *Finding {
	ruleFlag := false

	if _, ok := r.Keys[strings.ToLower(key)]; ok {
		ruleFlag = true
	}

	if ruleFlag {
		return &Finding{
			Severity:       r.Base.severity,
			Rule:           "Правило паролей",
			Path:           path,
			Message:        fmt.Sprintf("Передаются чувствительные данные %q", key),
			Recommendation: r.Base.recommendation,
		}
	}
	return nil
}

type PublicBindRule struct {
	Base            baseRule
	Keys            map[string]struct{}
	UnsafeAddresses map[string]struct{}
}

func NewPublicBindRule(cfgRule config.PublicBindRuleConfig) *PublicBindRule {
	keyMap := make(map[string]struct{})
	adrMap := make(map[string]struct{})

	for _, key := range cfgRule.Keys {
		keyMap[strings.ToLower(key)] = struct{}{}
	}

	for _, adr := range cfgRule.UnsafeAddresses {
		adrMap[strings.ToLower(adr)] = struct{}{}
	}

	return &PublicBindRule{
		Base:            newBaseRule(cfgRule.CommonRuleConfig),
		Keys:            keyMap,
		UnsafeAddresses: adrMap,
	}
}

func (r *PublicBindRule) Check(path string, key string, value any) *Finding {
	ruleFlag := false
	stringValue, ok := value.(string)
	if !ok {
		return nil
	}

	if _, ok := r.Keys[strings.ToLower(key)]; ok {
		if _, ok := r.UnsafeAddresses[strings.ToLower(stringValue)]; ok {
			ruleFlag = true
		}
	}

	if ruleFlag {
		return &Finding{
			Severity:       r.Base.severity,
			Rule:           "Правило публичных адресов",
			Path:           path,
			Message:        fmt.Sprintf("Публичный адрес %q:%q", key, stringValue),
			Recommendation: r.Base.recommendation,
		}
	}
	return nil
}

type TLSRule struct {
	Base  baseRule
	FKeys map[string]struct{}
	TKeys map[string]struct{}
}

func NewTLSRule(cfgRule config.DisabledTLSRuleConfig) *TLSRule {
	tKeyMap := make(map[string]struct{})
	fKeyMap := make(map[string]struct{})

	for _, key := range cfgRule.TrueKeys {
		tKeyMap[strings.ToLower(key)] = struct{}{}
	}

	for _, adr := range cfgRule.FalseKeys {
		fKeyMap[strings.ToLower(adr)] = struct{}{}
	}

	return &TLSRule{
		Base:  newBaseRule(cfgRule.CommonRuleConfig),
		FKeys: fKeyMap,
		TKeys: tKeyMap,
	}
}

func (r *TLSRule) Check(path string, key string, value any) *Finding {
	ruleFlag := false
	boolValue, ok := value.(bool)
	if !ok {
		return nil
	}

	if _, ok := r.FKeys[strings.ToLower(key)]; ok {
		if !boolValue {
			ruleFlag = true
		}
	}

	if _, ok := r.TKeys[strings.ToLower(key)]; ok {
		if boolValue {
			ruleFlag = true
		}
	}

	if ruleFlag {
		return &Finding{
			Severity:       r.Base.severity,
			Rule:           "Правило TLS",
			Path:           path,
			Message:        fmt.Sprintf("Некорректная настройка TLS %q:%v", key, boolValue),
			Recommendation: r.Base.recommendation,
		}
	}
	return nil
}

type AlgRule struct {
	Base             baseRule
	Keys             map[string]struct{}
	UnsafeAlgorithms map[string]struct{}
}

func NewAlgRule(cfgRule config.WeakAlgorithmRuleConfig) *AlgRule {
	keyMap := make(map[string]struct{})
	algMap := make(map[string]struct{})

	for _, key := range cfgRule.Keys {
		keyMap[strings.ToLower(key)] = struct{}{}
	}

	for _, alg := range cfgRule.UnsafeAlgorithms {
		algMap[strings.ToLower(alg)] = struct{}{}
	}

	return &AlgRule{
		Base:             newBaseRule(cfgRule.CommonRuleConfig),
		Keys:             keyMap,
		UnsafeAlgorithms: algMap,
	}
}

func (r *AlgRule) Check(path string, key string, value any) *Finding {
	ruleFlag := false
	stringValue, ok := value.(string)
	if !ok {
		return nil
	}

	if _, ok := r.Keys[strings.ToLower(key)]; ok {
		if _, ok := r.UnsafeAlgorithms[strings.ToLower(stringValue)]; ok {
			ruleFlag = true
		}
	}

	if ruleFlag {
		return &Finding{
			Severity:       r.Base.severity,
			Rule:           "Правило небезопасных алгоритмов",
			Path:           path,
			Message:        fmt.Sprintf("Слишком слабый алгоритм - %q", stringValue),
			Recommendation: r.Base.recommendation,
		}
	}
	return nil
}

// Build - фабрика списка правил
func Build(cfg config.RulesConfig) []Rule {
	rules := make([]Rule, 0)

	if cfg.DebugLogging.Enabled {
		r := NewLogRule(cfg.DebugLogging)
		rules = append(rules, r)
	}

	if cfg.PlaintextPassword.Enabled {
		r := NewPassRule(cfg.PlaintextPassword)
		rules = append(rules, r)
	}

	if cfg.PublicBind.Enabled {
		r := NewPublicBindRule(cfg.PublicBind)
		rules = append(rules, r)
	}

	if cfg.DisabledTLS.Enabled {
		r := NewTLSRule(cfg.DisabledTLS)
		rules = append(rules, r)
	}

	if cfg.WeakAlgorithm.Enabled {
		r := NewAlgRule(cfg.WeakAlgorithm)
		rules = append(rules, r)
	}

	return rules
}
