

# Imagen de Docker para WordPress con Xdebug

Esta imagen extiende la imagen oficial de WordPress con Xdebug instalado y configurado para depuración paso a paso en el puerto 9000. Aunque Xdebug 3 utiliza el puerto 9003 de forma predeterminada, esta imagen conserva el puerto 9000 para mantener la compatibilidad con los usuarios existentes.

## Uso

Configure Xdebug en tiempo de ejecución con `XDEBUG_MODE` y `XDEBUG_CONFIG`. Por ejemplo, establezca `discover_client_host=true` para que Xdebug se conecte al cliente que realizó la solicitud HTTP. Cuando la detección del cliente no sea adecuada, establezca explícitamente `client_host=<host>`.

Docker Desktop proporciona `host.docker.internal`. En Linux con Docker Engine, asigne ese nombre de host a `host-gateway` de Docker, como se muestra a continuación.

Consulte la [documentación de configuraciones de Xdebug](https://xdebug.org/docs/all_settings) para ver todas las opciones disponibles.

## Docker Compose

Ejemplo de `compose.yml`:

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

Instale la [extensión PHP Debug](https://marketplace.visualstudio.com/items?itemName=xdebug.php-debug) y configure `pathMappings` para que VS Code pueda asignar las rutas del contenedor a los archivos locales.

Ejemplo de `.vscode/launch.json`:

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

## Construcción

De forma predeterminada, la compilación local utiliza la última imagen de WordPress y la versión estable de Xdebug:

```sh
docker build -t wordpress-xdebug .
```

Para una combinación de versiones reproducible, pase ambos argumentos de compilación:

```sh
docker build \
  --build-arg WORDPRESS_VERSION=7.0.2 \
  --build-arg XDEBUG_VERSION=3.5.3 \
  -t wordpress-xdebug:wp7.0.2-xdebug3.5.3 \
  .
```
