cd ..
docker run --rm --mount type=bind,src=${PWD},dst=/src --workdir /src golang:alpine sh -c "cd data_fetcher && go build -o=./data-fetcher -trimpath -mod=readonly -ldflags=-s -ldflags=-w ."
cd data_fetcher
docker build -t data-fetcher .
rm data-fetcher
