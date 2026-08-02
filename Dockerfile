ARG WORDPRESS_VERSION=latest
FROM wordpress:${WORDPRESS_VERSION}

ARG XDEBUG_VERSION

LABEL maintainer="André Costa <andreccosta@me.com>"

ENV XDEBUG_PORT=9000
ENV XDEBUG_IDEKEY=docker

RUN pecl install "xdebug${XDEBUG_VERSION:+-${XDEBUG_VERSION}}" \
    && docker-php-ext-enable xdebug

RUN printf '%s\n' \
    'xdebug.mode=debug' \
    'xdebug.start_with_request=yes' \
    "xdebug.client_port=${XDEBUG_PORT}" \
    "xdebug.idekey=${XDEBUG_IDEKEY}" \
    > /usr/local/etc/php/conf.d/xdebug.ini
