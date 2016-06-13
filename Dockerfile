#FROM gliderlabs/alpine:latest
FROM index.shurenyun.com/zqdou/ubuntu-go:1.5.1
#VOLUME /mnt/routes
COPY . /src
RUN cd /src && ./build.sh
ENTRYPOINT ["/bin/start.sh"]
