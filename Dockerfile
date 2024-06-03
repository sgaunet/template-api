FROM alpine:3.20.0 AS builder
RUN apk add --no-cache ca-certificates

FROM scratch
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY resources /
COPY webserver               /opt/webserver/
# COPY static                             /opt/webserver/static/
# COPY templates                          /opt/webserver/templates
WORKDIR /opt/webserver
USER MyUser
EXPOSE 3000
CMD [ "/opt/webserver/webserver" ]
