# swagger-plugin

Для генерации документации в формате [openAPI](https://swagger.io/docs/specification/about/), необходимо выполнить
команду:

```bash
tg plugin run swagger-plugin --ServicePath . --OutFile ../api/swagger.yaml
```

Где,

`ServicePath` - путь до папки с интерфейсом (в норме для `tg` эта папка является рабочей)
`OutFile` - путь, где будет сохранён результат
