# WordPress with Xdebug Docker image

This image extends the official WordPress image with Xdebug installed and
configured for step debugging on port 9000. Although Xdebug 3 defaults to port
9003, this image retains port 9000 for compatibility with existing users.

## Usage

Configure Xdebug at runtime with `XDEBUG_MODE` and `XDEBUG_CONFIG`. For example,
set `discover_client_host=true` to have Xdebug connect to the client that made
the HTTP request. When client discovery is not suitable, set
`client_host=<host>` explicitly.

Docker Desktop provides `host.docker.internal`. On Linux with Docker Engine,
map that hostname to Docker's `host-gateway`, as shown below.

See the [Xdebug settings documentation](https://xdebug.org/docs/all_settings)
for all available options.

## Docker Compose

Example `compose.yml`:

```yaml
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
      XDEBUG_CONFIG: client_host=host.docker.internal client_port=9000
```

## Visual Studio Code

Install the [PHP Debug extension](https://marketplace.visualstudio.com/items?itemName=xdebug.php-debug)
and configure `pathMappings` so VS Code can map container paths to local files.

Example `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Listen for Xdebug",
      "type": "php",
      "request": "launch",
      "port": 9000,
      "pathMappings": {
        "/var/www/html": "${workspaceFolder}/wp"
      }
    }
  ]
}
```

## Building

By default, a local build uses the latest WordPress image and stable Xdebug
release:

```sh
docker build -t wordpress-xdebug .
```

For a reproducible version combination, pass both build arguments:

```sh
docker build \
  --build-arg WORDPRESS_VERSION=7.0.2 \
  --build-arg XDEBUG_VERSION=3.5.3 \
  -t wordpress-xdebug:wp7.0.2-xdebug3.5.3 \
  .
```
