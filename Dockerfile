FROM wordpress:latest
LABEL maintainer André Costa <andreccosta@me.com>

ENV XDEBUG_PORT 9000

RUN yes | pecl install xdebug && \
  echo "zend_extension=$(find /usr/local/lib/php/extensions/ -name xdebug.so)" > /usr/local/etc/php/conf.d/xdebug.ini && \
  echo "xdebug.remote_enable=on" >> /usr/local/etc/php/conf.d/xdebug.ini && \
  echo "xdebug.remote_autostart=on" >> /usr/local/etc/php/conf.d/xdebug.ini && \
  echo "xdebug.idekey=REMOTE" >> /usr/local/etc/php/conf.d/xdebug.ini

EXPOSE 9000

