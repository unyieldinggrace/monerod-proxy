# Docker How-To

The Docker deployment adds a Traefik balancer and adds SSL support making use of Let's Encrypt. By default, it also adds a CORS headers middleware, so the node can be used in web applications seamlessly.

### Requirements

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- A DNS domain pointing to your server (optional)

### Build and deploy

To build the base image, run:

```
docker-compose build
```

Once completed, configure the config file `.env` and `config.docker.ini` to your needs. Don't change the port and keep the node `http://monerod:18089/` as a primary to use the local node.

Then the node is ready to deploy with:

```
docker-compose pull
docker-compose up -d
```
