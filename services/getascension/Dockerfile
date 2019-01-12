FROM alpine:latest

RUN apk --no-cache add ca-certificates
# RUN apt-get update
# RUN apt-get install -y ca-certificates

ENV INSIDE_DOCKER 1

EXPOSE 443 80

RUN mkdir -p /root
ENV SOCIALDIR /root

# ADD executables/gops /usr/bin
ADD executables/server /getascension.com

CMD ["/getascension.com"]