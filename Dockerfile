FROM alpine:latest
RUN apk --no-cache add tzdata  && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone
ADD test_v1 /test_v1
ENTRYPOINT [ "/test_v1" ]