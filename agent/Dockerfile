FROM centos:latest
MAINTAINER Techie York "york528@yeah.net"
ENV REFRESHED_AT 2017-05-13
ENV LD_LIBRARY_PATH=/usr/local/lib

ADD ./admin/ /data/soft/DarkMetrix/monitor/agent/admin
ADD ./conf/ /data/soft/DarkMetrix/monitor/agent/conf
ADD ./bin/ /data/soft/DarkMetrix/monitor/agent/bin
ADD ./plugin/ /data/soft/DarkMetrix/monitor/agent/plugin
ADD ./lib/libstatgrab/centos/libstatgrab.so.10.0.0 /usr/local/lib

RUN ln -s /usr/local/lib/libstatgrab.so.10.0.0 /usr/local/lib/libstatgrab.so.10
RUN ln -s /usr/local/lib/libstatgrab.so.10.0.0 /usr/local/lib/libstatgrab.so

VOLUME ["/var/tmp"]
WORKDIR "/data/soft/DarkMetrix/monitor/agent/bin/"

ENTRYPOINT ["./dm_monitor_agent"]
CMD ["-h"]
