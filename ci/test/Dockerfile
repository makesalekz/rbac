FROM golang:latest AS builder

COPY . /src
WORKDIR /src
ARG TOKEN

RUN git config --global url.https://gitlab-ci-token:${TOKEN}@gitlab.calendaria.team.insteadOf https://gitlab.calendaria.team && \
    export GOPRIVATE=gitlab.calendaria.team

CMD ["go", "test", "-v", "-count=1", "./..."]