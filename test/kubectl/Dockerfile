FROM alpine

LABEL maintainer="Robert Buck <bob@continuul.io>"

COPY ./discover /usr/local/bin/discover

ENTRYPOINT [ "/usr/local/bin/discover" ]

CMD [ "addrs", "provider=kubernetes", "namespace=hashicorp", "label_key=role", "label_value=discoverable" ]
