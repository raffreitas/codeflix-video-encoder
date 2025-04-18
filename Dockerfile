FROM golang:1.24-alpine

ENV PATH="$PATH:/opt/bento4/bin"
ENV PYTHONPATH="/opt/Bento4/Source/Python"

RUN apk add ffmpeg bash make git python3 unzip gcc g++ cmake && \
    cd /opt && \
    git clone https://github.com/axiomatic-systems/Bento4.git && \
    mkdir -p Bento4/cmakebuild && \
    cd Bento4/cmakebuild && \
    cmake -DCMAKE_BUILD_TYPE=Release .. && \
    make -j$(nproc) && \
    mkdir -p /opt/bento4/bin && \
    find /opt/Bento4/cmakebuild -type f -perm +111 -name "mp4*" -exec cp {} /opt/bento4/bin/ \; && \
    find /opt/Bento4/Source/Python -type f -name "*.py" -exec cp {} /opt/bento4/bin/ \; && \
    chmod +x /opt/bento4/bin/*.py && \
    ln -s /opt/bento4/bin/mp4-dash.py /opt/bento4/bin/mp4dash && \
    apk del git cmake && \
    rm -rf /var/cache/apk/* /opt/Bento4/cmakebuild

WORKDIR /go/src

ENTRYPOINT ["tail", "-f", "/dev/null"]
