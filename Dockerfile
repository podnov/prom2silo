FROM centos:7

RUN useradd -M prom2silo

ADD prom2silo /

RUN chown prom2silo:prom2silo /prom2silo

USER prom2silo

CMD ["/prom2silo"]

