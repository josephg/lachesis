FROM glider as stage0

# Glider can be found in https://github.com/andrecronje/evm/blob/master/docker/glider/Dockerfile

RUN mkdir -p "$GOPATH/src/github.com/andrecronje/lachesis" /cp_bin /bin && \
    git clone 'https://github.com/SamuelMarks/lachesis' "$GOPATH/src/github.com/andrecronje/lachesis" && \
    cd "$GOPATH/src/github.com/andrecronje/lachesis" && \
    glide install && \
    cd "$GOPATH/src/github.com/andrecronje/lachesis/cmd/lachesis" && \
    go build -ldflags "-linkmode external -extldflags -static -s -w" -a main.go && \
    mv "$GOPATH/src/github.com/andrecronje/lachesis/cmd/lachesis/main" /cp_bin/lachesis

# ADD https://github.com/upx/upx/releases/download/v3.95/upx-3.95-amd64_linux.tar.xz /tmp
COPY upx /cp_bin/

RUN apk --no-cache add libc-dev cmake && \
    # /bin/tar --version && \
    # md5sum /tmp/upx-3.95-amd64_linux.tar.xz && \
    # /bin/tar xf /tmp/upx-3.95-amd64_linux.tar.xz && \
    # /bin/tar -C /bin --strip-components=1 -xzf /tmp/upx-3.95-amd64_linux.tar.xz && \
    git clone https://github.com/SamuelMarks/docker-static-bin /build/docker-static-bin -b boxed-main && \
    mkdir /build/docker-static-bin/c/cmake-build-release && \
    cd    /build/docker-static-bin/c/cmake-build-release && \
    cmake -DCMAKE_BUILD_TYPE=Release .. && make && \
    mv main /cp_bin/cp_bin && cd /cp_bin \
    ln -s /cp_bin/cp_bin /cp_bin/copy && \
    ln -s /cp_bin/cp_bin /cp_bin/crappy_sh && \
    ln -s /cp_bin/cp_bin /cp_bin/env && \
    ln -s /cp_bin/cp_bin /cp_bin/list && \
    ln -s /cp_bin/cp_bin /cp_bin/symlink && \
    strip -s /cp_bin/cp_bin /cp_bin/lachesis && \
    /cp_bin/upx --brute /cp_bin/lachesis /cp_bin/cp_bin

FROM alpine as lachesis_base

ENV node_num=0
ENV node_addr='127.0.0.1'

EXPOSE 1338
EXPOSE 1339
EXPOSE 8000
EXPOSE 12000

# cp -r /etc/ssl/certs certs, then add to your `docker build`: `--build-arg ca_certificates=certs`
ARG ca_certificates
ADD "$ca_certificates" /etc/ssl/certs/
COPY --from=0 /cp_bin /cp_bin

COPY peers.json /lachesis_data_dir/
COPY nodes /nodes

ENTRYPOINT ["/cp_bin/crappy_sh", "-v", "-e", "-c", "/cp_bin/upx -d /cp_bin/cp_bin /cp_bin/lachesis ; /cp_bin/env ; /cp_bin/list /cp_bin ; /cp_bin/list /nodes ; /cp_bin/copy /nodes/$node_num/priv_key.pem /lachesis_data_dir/priv_key.pem ; /cp_bin/list /lachesis_data_dir ; /cp_bin/lachesis run --datadir /lachesis_data_dir --store_path /lachesis_data_dir/badger_db -node_addr=$node_addr:12000 -heartbeat=100 -no_client"]
