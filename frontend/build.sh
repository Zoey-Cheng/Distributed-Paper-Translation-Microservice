npm run build
docker build . -t paper-translation-web:$(git rev-parse --short HEAD)
echo TAG=$(git rev-parse --short HEAD) > .env
