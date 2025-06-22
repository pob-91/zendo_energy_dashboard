rm -r node_modules

docker run --rm --mount type=bind,src=${PWD},dst=/src --workdir /src node:jod-alpine sh -c "npm install && npm run build"

docker build -t zendo-web-app .

rm -r dist
rm -r node_modules

npm install
