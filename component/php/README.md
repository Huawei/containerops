## Build PHP Component

1.build PHP base image

```bash
cd component/iamges/php
docker build -t containerops/php:7.1.4 --build-arg php_version=7.1.4  .
```

2.build composer base image

```bash
cd component/php/composer
docker build -t containerops/composer:latest .
```

3.build component image

```bash
cd component/php/component-composer
docker build -t containerops/component-composer:latest .
```

4.run component

```bash
docker run --env CO_DATA="git-url=https://github.com/TIGERB/easy-php.git action=install" containerops/component-composer:latest
```