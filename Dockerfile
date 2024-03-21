FROM scratch
COPY lll-fixer /usr/bin/lll-fixer
ENTRYPOINT ["/usr/bin/lll-fixer"]
CMD ["help"]
