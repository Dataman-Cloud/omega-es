FROM index.shurenyun.com/zqdou/ubuntu-go:1.5.1
COPY . /src
RUN cd /src && ./build.sh
ENTRYPOINT ["/bin/start.sh"]
