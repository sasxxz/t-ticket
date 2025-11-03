FROM alpine:latest
WORKDIR /Ticket
RUN mkdir ./server
ADD ./ticket ./
COPY ./web ./server/web
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai
RUN ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime
EXPOSE 1021
ENTRYPOINT ["/Ticket/ticket"]
