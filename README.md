## RUN

The `config_example.yaml` could be used as `config.yaml`.

```shell
 cp ./config/config_example.yaml ./config/config.yaml
```

### Docker

Requires `config.yaml`.
```shell
docker-compose up -d
```

### Local

Requires `config.yaml`.
```shell
go run ./cmd/main.go -c <full_path_to_config>/config.yaml
```
