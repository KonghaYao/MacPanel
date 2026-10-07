export const DEFAULT_COMPOSE_TEMPLATE = `services:
  app:
    image: nginx:latest
    container_name: my-app
    restart: unless-stopped
    ports:
      - "8080:80"
    volumes:
      - ./data:/usr/share/nginx/html
    networks:
      - 1panel-network

networks:
  1panel-network:
    external: true
`;

export const DEFAULT_COMPOSE_TEMPLATE_ID = 0;
