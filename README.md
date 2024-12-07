# Simple Collaborative Chat CLI with NATS and Go 🗨️

Este proyecto es una aplicación CLI básica de chat creada en Go, que se conecta a un servidor NATS para enviar y recibir mensajes de chat en tiempo real.

## Objetivo 🎯

Crear una aplicación de chat sencilla en la línea de comandos utilizando Go, que interactúe con un servidor NATS y permita a los usuarios enviar y recibir mensajes en un canal determinado.

## Requisitos 📋

Antes de ejecutar la aplicación, asegúrate de tener instalados los siguientes componentes:

- [Go 1.18+](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started) (para ejecutar el servidor NATS)

## Instalación 🛠️

1. Clona este repositorio:

   ```bash
   git clone https://github.com/brivaro/Chat-NATS
   cd Chat-NATS
   ```

2. Ejecuta el contenedor de Docker para levantar NATS:

   ```bash
   docker compose up --build
   ```

   Esto iniciará un servidor NATS con JetStream habilitado, accesible en `nats://localhost:4222`.


## Uso 🚀

1. Ejecuta la aplicación en tu terminal:

   ```bash
   go run main.go
   ```

2. Te pedirá que ingreses los siguientes datos:

   - **Server URL**: La URL del servidor NATS (por defecto `nats://localhost:4222`).
   - **Channel**: El nombre del canal al que te unirás (Formato de los chats: "chat.>", por ejemplo, chat.brian, chat.brian.hoy.se.sale.de.fiesta...).
   - **User**: Tu nombre de usuario (sin espacios ni caracteres especiales como indica el mss).

   Ejemplo:

   ```
   Please, insert the next values to join a chat:
   Server URL: nats://localhost:4222
   Channel (Format: chat.>): chat.brian
   User (Without: whitespace, ., *, >, path separators (forward or backward slash), or non-printable characters): alice
   ```

3. Una vez que hayas ingresado los datos, se suscribirá al canal y podrás comenzar a enviar y recibir mensajes en tiempo real. Si te sales del chat, y vuelves a ejecutar la app, `los mss de la hora pasada persisten`.

4. Para salir, presiona `Ctrl+C`. La aplicación limpiará los recursos antes de cerrar.

## Docker Compose 🔧

Este proyecto utiliza Docker Compose para levantar el servidor NATS. El archivo `docker-compose.yml` está configurado de la siguiente manera:

```yaml
version: '3'

services:
  nats:
    image: nats:latest
    container_name: nats
    command: --js
    ports:
    - 4222:4222
    volumes:
      - chat_data:/tmp/nats/jetstream

volumes:
  chat_data:
```

Cuenta con un volumen persistente, es decir, una vez se cierren o borren los contenedores, si volvemos a levantar todo, los mensajes persistirán.

## Contribuciones 🤝

Si tienes sugerencias o mejoras para el proyecto, ¡no dudes en hacer un pull request! 

---

¡Disfruta chateando! 🎉

---