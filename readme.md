# Wordpress XDebug Docker image

## Usage

The environment variable **XDEBUG_CONFIG** let's you configure the XDebug PHP extension. For example to provide the Docker host IP address set it as `remote_host=<host ip>`. For Docker 18.03.x and up you should just be able to use `remote_host=host.docker.internal`.

## Docker Compose

Example configuration file `docker-compose.yml`:

```yml
version: '3.3'

services:
  db:
    image: mysql
    restart: on-failure
    environment:
      MYSQL_ROOT_PASSWORD: somewordpress
      MYSQL_DATABASE: wordpress
      MYSQL_USER: wordpress
      MYSQL_PASSWORD: wordpress

  wp:
    depends_on:
    - db
    image: andreccosta/wordpress-xdebug
    volumes:
    - ./wp:/var/www/html
    ports:
    - 8080:80
    restart: on-failure
    environment:
      WORDPRESS_DB_HOST: db:3306
      WORDPRESS_DB_USER: wordpress
      WORDPRESS_DB_PASSWORD: wordpress
      XDEBUG_CONFIG: remote_host=host.docker.internal
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
      "port": 9000,
      "pathMappings": {
        "/var/www/html": "${workspaceRoot}/wp",
      }
    }
  ]
}
```
