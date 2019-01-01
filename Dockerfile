FROM alpine:3.7
RUN apk add --no-cache bash curl bind-tools tcpdump
RUN mkdir -p /usr/bin
ADD bin/ /usr/bin/