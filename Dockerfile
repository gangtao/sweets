FROM alpine:3.7
LABEL author="gang.tao@outlook.com"

RUN mkdir /home/sweets
WORKDIR /home/sweets

COPY ./server/server /home/sweets/
COPY ./client/client /home/sweets/

EXPOSE 8080

ENTRYPOINT ["/home/sweets/server"]