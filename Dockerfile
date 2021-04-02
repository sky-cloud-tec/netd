FROM registry.cn-hangzhou.aliyuncs.com/sky-cloud-tec/alpine:alpine
MAINTAINER songtianyi@sky-cloud.net

ADD netd /usr/bin/
ADD cfg.ini /etc/netd/cfg.ini
RUN mkdir -p /var/log/netd


RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
		  echo "Asia/Shanghai" > /etc/timezone

CMD /usr/bin/netd --cfg /etc/netd/cfg.ini
