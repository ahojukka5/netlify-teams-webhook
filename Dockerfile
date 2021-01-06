FROM golang AS app
COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build -a -ldflags '-s' .

FROM alpine AS certs
# https://medium.com/the-go-journey/x509-certificate-signed-by-unknown-authority-running-a-go-app-inside-a-docker-container-a12869337eb
RUN apk --no-cache add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=app /build/netlify-teams-webhook .
ENTRYPOINT ["./netlify-teams-webhook"]
