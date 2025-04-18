FROM golang:1.24-alpine

ENV PATH="$PATH:/bin/bash" \
    BENTO4_BIN="/opt/bento4/bin" \
    PATH="$PATH:/opt/bento4/bin"

RUN apk add --no-cache ffmpeg bash make wget python3 unzip gcc g++ scons

WORKDIR /tmp/bento4
ENV BENTO4_PATH="/opt/bento4"

RUN wget -q "https://www.bok.net/Bento4/binaries/Bento4-SDK-1-6-0-640.x86_64-unknown-linux.zip" && \
    mkdir -p ${BENTO4_PATH} && \
    unzip Bento4-SDK-1-6-0-640.x86_64-unknown-linux.zip -d /tmp && \
    mv /tmp/Bento4-SDK-1-6-0-640.x86_64-unknown-linux/* ${BENTO4_PATH}/ && \
    rm -rf Bento4-SDK-1-6-0-640.x86_64-unknown-linux.zip /tmp/Bento4-SDK-1-6-0-640.x86_64-unknown-linux && \
    apk del wget unzip && \
    rm -rf /var/cache/apk/* /tmp/*

WORKDIR /go/src

ENTRYPOINT ["tail", "-f", "/dev/null"]