FROM busybox
LABEL maintainer "Qiniu ATLab <ai@qiniu.com>"

COPY dbstorage /workspace/dbstorage
COPY dbstorage.conf /workspace/dbstorage.conf
WORKDIR /workspace/
CMD [ "./dbstorage", "-f", "./dbstorage.conf"]
