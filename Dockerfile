FROM hub.sky-cloud.net/cicd/alpine:alpine
MAINTAINER songtianyi@sky-cloud.net

ADD netd /usr/bin/

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
		  echo "Asia/Shanghai" > /etc/timezone

CMD /usr/bin/netd --loglevel ${LOG_LEVEL} --logfile ${LOG_FILE} \
		jrpc --address ${ADDR}
