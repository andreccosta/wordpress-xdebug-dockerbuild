# Wordpress with XDebug Docker image

## Usage

The XDebug extension can be configured via environment variables. Namely **XDEBUG_MODE** and **XDEBUG_CONFIG**.

For example to allow XDebug to try to automatically connect back to the client that made the HTTP request you would add `discover_client_host=true` to **XDEBUG_CONFIG**. Or in scenarios where that is not feasible you would provide the Docker host IP address and set it as `client_host=<host ip>`.

For Docker 18.03.x and up, non-Linux users should be able to just use `client_host=host.docker.internal`. Linux users should configure `extra_hosts` with `host-gateway` instead.

You can check additional information about what XDebug settings are available in the documentation [here](https://xdebug.org/docs/all_settings).

## Docker Compose

Example configuration file `docker-compose.yml`:

```yml
services:
  db:
    image: mariadb:11
    restart: on-failure
    environment:
      MARIADB_ROOT_PASSWORD: wordpress
      MARIADB_DATABASE: wordpress
      MARIADB_USER: wordpress
      MARIADB_PASSWORD: wordpress

  wp:
    depends_on:
      - db
    image: andreccosta/wordpress-xdebug:latest
    volumes:
      - ./wp:/var/www/html
    ports:
      - 8080:80
    restart: on-failure
    extra_hosts:
      - "host.docker.internal:host-gateway"
    environment:
      WORDPRESS_DB_HOST: db:3306
      WORDPRESS_DB_USER: wordpress
      WORDPRESS_DB_PASSWORD: wordpress
      XDEBUG_MODE: debug
      XDEBUG_CONFIG: start_with_request=yes client_host=host.docker.internal client_port=9003
```

## Visual Studio Code

To use XDebug in Visual Studio Code you need the [PHP Debug extension](https://marketplace.visualstudio.com/items?itemName=felixfbecker.php-debug).

Also to make VS Code map the paths on the container to the ones on the host, you have to set the pathMappings settings in your `launch.json`.

Example configuration file `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Listen for XDebug",
      "type": "php",
      "request": "launch",
      "port": 9003,
      "pathMappings": {
        "/var/www/html": "${workspaceRoot}/wp",
      }
    }
  ]
}
```
