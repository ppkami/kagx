FROM alpine:3.8

ENV KAGX_VERSION 0.0.1.2

RUN wget http://petlludhz.bkt.clouddn.com/kagx_v${KAGX_VERSION}_linux_amd64.tar.gz
RUN tar xzf kagx_v${KAGX_VERSION}_linux_amd64.tar.gz &&\
    mv kagx /usr/local/ &&\
    rm kagx_v${KAGX_VERSION}_linux_amd64.tar.gz

VOLUME /usr/local/kagx/conf

WORKDIR /usr/local/kagx

EXPOSE 40000/udp

ENTRYPOINT ["/usr/local/kagx/bin/kagxs"]
