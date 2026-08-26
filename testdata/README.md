# Проверки

| Сценарий                         | JSON | YAML | Ожидаемый результат |
|----------------------------------| --- | --- |---------------------|
| Safe configuration               | `safe.json` | `safe.yaml` | 0 находок           |
| All required rules               | `all-unsafe.json` | `all-unsafe.yaml` | 5 находок           |
| Nested objects and arrays        | `nested-arrays.json` | `nested-arrays.yaml` | 6 находок           |
| Alternative keys and letter case | `alternative-keys.json` | `alternative-keys.yaml` | 5 находок           |
| Safe and unsafe values together  | `mixed.json` | `mixed.yaml` | 3 находок           |
| Values of unsupported types      | `wrong-types.json` | `wrong-types.yaml` | 0 находок           |
| Invalid syntax                   | `invalid-syntax.json` | `invalid-syntax.yaml` | ошибка парсинга     |
| Non-object root                  | `invalid-root.json` | `invalid-root.yaml` | ошибка валидации    |
