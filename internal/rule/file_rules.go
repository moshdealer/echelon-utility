package rule

import (
	"fmt"
	"io/fs"

	"echelon-utility/internal/config"
)

const unsafeWritePermissions fs.FileMode = 0o022

// FileRule проверяет свойства самого конфигурационного файла (бонусное задание на os.Stat)
type FileRule interface {
	Check(path string, mode fs.FileMode) *Finding
}

type FilePermissionRule struct {
	Base baseRule
}

func NewFilePermissionRule(cfgRule config.FilePermissionRuleConfig) *FilePermissionRule {
	return &FilePermissionRule{
		Base: newBaseRule(cfgRule.CommonRuleConfig),
	}
}

func (r *FilePermissionRule) Check(path string, mode fs.FileMode) *Finding {
	permissions := mode.Perm()
	if permissions&unsafeWritePermissions == 0 {
		return nil
	}

	return &Finding{
		Severity:       r.Base.severity,
		Rule:           "Правило прав доступа к файлу",
		Path:           path,
		Message:        fmt.Sprintf("Небезопасные права доступа к файлу %04o", permissions),
		Recommendation: r.Base.recommendation,
	}
}

func BuildFileRules(cfg config.RulesConfig) []FileRule {
	rules := make([]FileRule, 0)

	if cfg.FilePermissions.Enabled {
		rules = append(rules, NewFilePermissionRule(cfg.FilePermissions))
	}

	return rules
}
