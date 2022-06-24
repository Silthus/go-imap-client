FROM scratch
COPY go-imap-client /
ENTRYPOINT ["/go-imap-client"]