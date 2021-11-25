FROM golang:1.17.3-buster as build_env 

WORKDIR /root/go_build

COPY . .

RUN make build

FROM  debian:buster

RUN   DEBIAN_FRONTEND=noninteractive \
      apt clean && \
      apt update && \
      apt install -y libreoffice && \
      apt clean

WORKDIR /var/www/html/file-server

COPY --from=build_env /root/go_build/app .
COPY --from=build_env /root/go_build/templates ./templates/
COPY --from=build_env /root/go_build/simsun.ttc /usr/share/fonts/

RUN chmod 644 /usr/share/fonts/simsun.ttc && fc-cache -fv

EXPOSE 8080

CMD ["/var/www/html/file-server/app"]
