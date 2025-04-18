FROM golang:1.24-alpine

RUN apk add --no-cache ffmpeg bash make wget git python3 unzip gcc g++ scons cmake

WORKDIR /tmp
RUN git clone https://github.com/axiomatic-systems/Bento4.git && \
    cd Bento4 && \
    mkdir cmakebuild && \
    cd cmakebuild && \
    cmake -DCMAKE_BUILD_TYPE=Release .. && \
    make && \
    mkdir -p /opt/bento4/bin && \
    find . -type f -perm +111 -name "mp4*" -exec cp {} /opt/bento4/bin/ \; && \
    cd /tmp && \
    rm -rf Bento4 && \
    apk del git cmake && \
    rm -rf /var/cache/apk/* /tmp/*

WORKDIR /go/src

ENTRYPOINT ["tail", "-f", "/dev/null"]