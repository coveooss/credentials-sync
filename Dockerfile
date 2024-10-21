FROM alpine@sha256:beefdbd8a1da6d2915566fde36db9db0b524eb737fc57cd1367effd16dc0d06d

# ca-certificates is needed to download files from S3
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# Install the program
ADD credentials-sync /app/credentials-sync
WORKDIR /app

ENTRYPOINT ["/app/credentials-sync"]