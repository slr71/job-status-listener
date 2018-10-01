FROM golang:1.11-alpine

ARG git_commit=unknown
ARG version="2.9.0"
ARG descriptive_version=unknown

COPY . /go/src/github.com/cyverse-de/job-status-listener
ENV CGO_ENABLED=0
RUN go install -v -ldflags "-X main.appver=$version -X main.gitref=$git_commit" github.com/cyverse-de/job-status-listener

EXPOSE 60000
ENTRYPOINT ["job-status-listener"]
CMD ["--help"]

LABEL org.label-schema.vcs-url="https://github.com/cyverse-de/job-status-listener"
LABEL org.label-schema.vcs-ref="$git_commit"
LABEL org.label-schema.version="$descriptive_version"
LABEL org.cyverse.git-ref="$git_commit"
LABEL org.cyverse.version="$version"
LABEL org.cyverse.descriptive-version="$descriptive_version"
