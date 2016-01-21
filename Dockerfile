FROM gliderlabs/alpine:latest
VOLUME /mnt/routes
COPY . /src
RUN cd /src && ./build.sh
ENTRYPOINT ["/bin/start.sh"]
