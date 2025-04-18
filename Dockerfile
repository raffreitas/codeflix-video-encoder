FROM golang:1.24-alpine

ENV PATH="$PATH:/opt/bento4/bin"

RUN apk add --no-cache \
    ffmpeg \
    bash \
    make \
    wget \
    git \
    python3 \
    unzip \
    gcc \
    g++ \
    scons \
    cmake

WORKDIR /opt
RUN git clone https://github.com/axiomatic-systems/Bento4.git && \
    mkdir -p Bento4/cmakebuild && \
    cd Bento4/cmakebuild && \
    cmake -DCMAKE_BUILD_TYPE=Release .. && \
    make -j$(nproc)

RUN mkdir -p /opt/bento4/bin

RUN find /opt/Bento4/cmakebuild -type f -perm +111 -name "mp4*" -exec cp {} /opt/bento4/bin/ \;

RUN find /opt/Bento4/Source/Python -type f -name "*.py" -exec cp {} /opt/bento4/bin/ \; && \
    chmod +x /opt/bento4/bin/*.py && \
    ln -s /opt/bento4/bin/mp4-dash.py /opt/bento4/bin/mp4dash

ENV PYTHONPATH="/opt/Bento4/Source/Python"

RUN apk del git cmake && \
    rm -rf /var/cache/apk/*

WORKDIR /go/src

ENTRYPOINT ["tail", "-f", "/dev/null"]
