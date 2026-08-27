package rule

// Rule - интерфейс, отвечающий за проверку строчки на соответствие заявленным правилам
type Rule interface {
	Check(path string, key string, value any) *Finding
}
