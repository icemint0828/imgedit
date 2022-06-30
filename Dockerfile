FROM scratch
ADD imgedit /usr/bin/imgedit
ENTRYPOINT ["/mnt/imgedit"]
