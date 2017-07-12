## Build PHP Component

1.build PHP base image

```bash
cd component/iamges/php
docker build -t containerops/php:7.1.4 --build-arg php_version=7.1.4  .
```

2.build composer base image

```bash
docker build -t containerops/composer:latest -f Base/composer/Dockerfile .
```

3.build component image

```bash
docker build -t containerops/component-composer:latest -f Dependence/component-composer/Dockerfile  .
```

4.run component

```bash
docker run --env CO_DATA="git-url=https://github.com/TIGERB/easy-php.git action=install" containerops/component-composer:latest
```