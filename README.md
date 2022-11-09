# go-musthave-shortener-tpl
Шаблон репозитория для практического трека «Go в веб-разработке».

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` - адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона выполните следующую команду:

```shell
git remote add -m main template https://github.com/yandex-praktikum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```shell
git fetch template && git checkout template/main .github
```

## swagger

```shell
swag init --output ./swagger/ -d cmd/shortener -g ../../internal/handlers/router.go --pd --parseInternal
```

```shell
swag fmt -d internal
```

## Go doc
```shell
godoc -http=:8080 -play
```

Затем добавьте полученные изменения в свой репозиторий.
