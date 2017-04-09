docker build -t tiki_ink /home/neo/docker/images/tiki_ink

# Stop it
docker rm -f tiki_ink_master
docker rm -f tiki_ink_worker0
docker rm -f tiki_ink_worker1


# Run a new master
docker run -d -v /home/neo/docker/data/jay/wwwroot/tiki/public/store:/tiki -p 30006:8080 --name tiki_ink_master tiki_ink /go/bin/ink -m

# Run two workers
docker run -d -v /home/neo/docker/data/jay/wwwroot/tiki/public/store:/tiki --name tiki_ink_worker0 tiki_ink
docker run -d -v /home/neo/docker/data/jay/wwwroot/tiki/public/store:/tiki --name tiki_ink_worker1 tiki_ink
