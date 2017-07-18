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
docker build -t containerops/phpmetrics:latest -f Analysis/phpmetrics/Dockerfile  .
docker build -t containerops/phploc:latest -f Analysis/phploc/Dockerfile  .
docker build -t containerops/phpcpd:latest -f Analysis/phpcpd/Dockerfile  .
docker build -t containerops/phpmd:latest -f Analysis/phpmd/Dockerfile  .
docker build -t containerops/phpcs:latest -f Analysis/phpcs/Dockerfile  .
```

4.run component

```bash
docker run --env CO_DATA="git-url=https://github.com/TIGERB/easy-php.git action=install" containerops/component-composer:latest 
docker run --env CO_DATA="git-url=http://192.168.123.201/yangkghjh/easy-php.git" containerops/phpmetrics:latest
docker run --env CO_DATA="git-url=http://192.168.123.201/yangkghjh/easy-php.git exclude=public" containerops/phploc:latest
docker run --env CO_DATA="git-url=http://192.168.123.201/yangkghjh/easy-php.git" containerops/phpcpd:latest
docker run --env CO_DATA="git-url=http://192.168.123.201/yangkghjh/easy-php.git" containerops/phpmd:latest
docker run --env CO_DATA="git-url=https://github.com/squizlabs/PHP_CodeSniffer.git report=full standard=phpcs.xml.dist" containerops/phpcs:latest
```