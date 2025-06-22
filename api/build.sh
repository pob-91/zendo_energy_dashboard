cd ..
docker run --rm --mount type=bind,src=${PWD},dst=/src --workdir /src golang:alpine sh -c "cd api && go build -o=./zendo-api -trimpath -mod=readonly -ldflags=-s -ldflags=-w ."
cd api
docker build -t zendo-api .
rm zendo-api
