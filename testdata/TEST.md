# Проверки

| Сценарий                                        | JSON                                          | YAML | Ожидаемый результат                | Фактический результат |
|-------------------------------------------------|-----------------------------------------------| --- |------------------------------------|-----------------------|
| Безопасный конфиг                               | `safe.json`                                   | `safe.yaml` | 0 находок                          | ОК                    |
| Нарушены все правила                            | `all-unsafe.json`                             | `all-unsafe.yaml` | 5 находок                          | ОК                    |
| Сложный конфиг с ошибками (вложенные структуры) | `nested-arrays.json`                          | `nested-arrays.yaml` | 6 находок                          | ОК                    |
| Поля в файле написаны CAPS'ом                   | `alternative-keys.json`                       | `alternative-keys.yaml` | 5 находок                          | ОК                    |
| Конфиг, где нарушены не все правила             | `mixed.json`                                  | `mixed.yaml` | 3 находок                          | ОК                    |
| Ошибки синтаксиса в конфиге                     | `invalid-syntax.json`                         | `invalid-syntax.yaml` | ошибка парсинга                    | ОК                    |
| Неверная структура конфига                      | `invalid-root.json`                           | `invalid-root.yaml` | ошибка валидации                   | ОК                    |
| Stdin                                           | `stdin-json.txt`                              | `stdin-yaml.txt` | 5 находок                          | ОК                    |
| Небезопасные права на файл                      | `file-permissions.json`                       | `file-permissions.yaml` | 1 находка уровня `MEDIUM`          | ОК                    |
| Обход директории                                | `testdata`                                    | `testdata` | Все ошибочные файлы попали в отчет | ОК                    |
| REST API                                        | `all-unsafe.json` через `POST /scan`          | `all-unsafe.yaml` через `POST /scan` | HTTP `200`, 5 находок              | ОК                    |
| gRPC API                                        | `all-unsafe.json` через `ScannerService/Scan` | `all-unsafe.yaml` через `ScannerService/Scan` | статус `OK`, 5 находок             | ОК                    |


# Команды для теста

Запускать следует из корня директории (папка echelon-utility)

## Проверка CLI
### Запуск cli для файлов
```
go run ./cli/cmd/echelon-utility testdata/all-unsafe.json
go run ./cli/cmd/echelon-utility testdata/all-unsafe.yaml
```
### Запуск cli для stdin
```
go run ./cli/cmd/echelon-utility --stdin
```
Вставляем текст конфигурации и нажимаем Ctrl+D

Либо
```
go run ./cli/cmd/echelon-utility --stdin < testdata/stdin-json.txt
go run ./cli/cmd/echelon-utility --stdin < testdata/stdin-yaml.txt
```

### Запуск cli с флагами
```
go run ./cli/cmd/echelon-utility -s testdata/all-unsafe.json
go run ./cli/cmd/echelon-utility --silent testdata/all-unsafe.json
```

### Запуск обхода директории
```
go run ./cli/cmd/echelon-utility testdata
go run ./cli/cmd/echelon-utility -s testdata
```
Примечание: может показаться, что silent не работает, но это не так - выдает ошибку из-за того что есть проблемные файлы invalid-syntax.* и invalid-root.* в каталоге (ошибка парсинга)

### Проверка на права доступа
Перед началом на всякий случай
```
chmod 666 testdata/file-permissions.json testdata/file-permissions.yaml
```
Сама проверка
```
go run ./cli/cmd/echelon-utility testdata/file-permissions.json
go run ./cli/cmd/echelon-utility testdata/file-permissions.yaml
```

## Проверка REST
```
docker-compose up --build -d
```

### Json
```
curl -sS \
  -X POST http://localhost:8080/scan \
  -H 'Content-Type: application/json' \
  --data-binary @testdata/all-unsafe.json | jq .
```
### Yaml
```
curl -sS \
  -X POST http://localhost:8080/scan \
  -H 'Content-Type: application/yaml' \
  --data-binary @testdata/all-unsafe.yaml | jq .
```

## Проверка gRPC

### Json
```
jq -Rs \
  '{sourceName:"all-unsafe.json", content:(@base64)}' \
  testdata/all-unsafe.json |
  grpcurl -plaintext -d @ \
  localhost:9090 \
  echelon.v1.ScannerService/Scan
```

### Yaml
```
jq -Rs \
  '{sourceName:"all-unsafe.yaml", content:(@base64)}' \
  testdata/all-unsafe.yaml |
  grpcurl -plaintext -d @ \
  localhost:9090 \
  echelon.v1.ScannerService/Scan
```

### Конец проверок
```
docker-compose down -v
```