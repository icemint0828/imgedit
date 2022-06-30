FROM scratch
ADD imgedit /usr/bin/imgedit
WORKDIR /mnt
ENTRYPOINT ["imgedit"]
