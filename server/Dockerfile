FROM wiremock/wiremock:2.32.0-alpine

WORKDIR /home/wiremock

COPY mappings ./mappings

ENTRYPOINT ["/docker-entrypoint.sh", "--global-response-templating"]
