FROM reg.qiniu.com/avaprd/argus-base-serving-base:20180719-v113-dev
LABEL maintainer "Qiniu ATLab <ai@qiniu.com>"

COPY dbstorage_tool /workspace/dbstorage_tool
COPY dbstorage_tool.conf /workspace/dbstorage_tool.conf
WORKDIR /workspace/
CMD [ "./dbstorage_tool", "-f", "./dbstorage_tool.conf"]
