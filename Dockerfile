FROM hub.sky-cloud.net/cicd/alpine:alpine
MAINTAINER songtianyi@sky-cloud.net

ENV LOG_LEVEL DEBUG
ENV LOG_FILE /var/log/netd/netd.log
ENV CONFIDENCE 30
ENV LOG_CFG_FLAG 0
ENV ADDR 0.0.0.0:8188

ADD netd /usr/bin/
RUN mkdir -p /var/log/netd


RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
		  echo "Asia/Shanghai" > /etc/timezone

CMD /usr/bin/netd --loglevel ${LOG_LEVEL} --logfile ${LOG_FILE} \
		jrpc --address ${ADDR} --confidence ${CONFIDENCE} --log-cfg-flag ${LOG_CFG_FLAG}
