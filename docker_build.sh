docker build -f server/Dockerfile -t qibla-backend-chat:latest . --no-cache

docker-compose down
docker-compose up -d --build
