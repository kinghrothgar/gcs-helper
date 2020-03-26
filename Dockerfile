FROM alpine:3.11.3
RUN  apk add --no-cache ca-certificates
ADD  gcs-helper /usr/bin/gcs-helper
RUN  adduser -u 999 -D gcs-helper
ENTRYPOINT ["/usr/bin/gcs-helper"]
