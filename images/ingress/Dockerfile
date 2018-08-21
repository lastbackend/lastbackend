FROM alpine:3.4

MAINTAINER lastbackend "team@lastbackend.com"

COPY ./images/ingress/build.sh /tmp/build.sh
RUN /tmp/build.sh

# forward request and error logs to docker log collector
RUN ln -sf /dev/stdout /var/log/nginx/access.log \
RUN ln -sf /dev/stderr /var/log/nginx/error.log

COPY ./images/ingress/rootfs/etc/nginx/lua /etc/nginx/lua
COPY ./images/ingress/rootfs/etc/nginx/nginx.conf /etc/nginx/nginx.conf
COPY ./images/ingress/rootfs/etc/nginx/conf.d /etc/nginx/conf.d
COPY ./images/ingress/rootfs/etc/nginx/html /etc/nginx/html

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]

