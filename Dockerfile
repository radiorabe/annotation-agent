FROM centos:8 AS build

ARG AUDIOWAVEFORM_VERSION=1.4.2

RUN dnf install -y \
    epel-release \
 && dnf install -y \
    --enablerepo PowerTools \
    boost-devel \
    cmake \
    gcc-c++ \
    gd-devel \
    libid3tag-devel \
    libmad-devel \
    libsndfile-devel \
    make \
 && dnf clean all

 RUN curl -s -L -o /tmp/audiowaveform.tar.gz https://github.com/bbc/audiowaveform/archive/$AUDIOWAVEFORM_VERSION.tar.gz \
  && tar xf /tmp/audiowaveform.tar.gz -C /tmp \
  && cd /tmp/audiowaveform-$AUDIOWAVEFORM_VERSION \
  && mkdir build \
  && cd build \
  && cmake -D ENABLE_TESTS=0 .. \
  && make \
  && make install

FROM centos:8

COPY --from=build /usr/local/bin/audiowaveform /usr/bin/

RUN dnf install -y \
    epel-release \
 && dnf install -y \
    --enablerepo PowerTools \
    boost-filesystem \
    boost-program-options \
    boost-regex \
    gd \
    libid3tag \
    libmad \
    libsndfile \
 && dnf clean all

ARG SONIC_ANNOTATOR_VERSION=1.6

RUN curl -s -o /tmp/sonic-annotator.tar.gz \
    https://code.soundsoftware.ac.uk/attachments/download/2708/sonic-annotator-$SONIC_ANNOTATOR_VERSION-linux64-static.tar.gz \
 && tar xf /tmp/sonic-annotator.tar.gz -C /tmp/ sonic-annotator-$SONIC_ANNOTATOR_VERSION-linux64-static/sonic-annotator \
 && mv /tmp/sonic-annotator-$SONIC_ANNOTATOR_VERSION-linux64-static/sonic-annotator /tmp/ \
 && cd /opt/ \
 && /tmp/sonic-annotator --appimage-extract \
 && mv squashfs-root sonic-annotator \
 && rm -rf \
    sonic-annotator/usr/lib/libtasn1.so.6 \
    /tmp/sonic-annotator.tar.gz \
    /tmp/sonic-annotator-$SONIC_ANNOTATOR_VERSION-linux64-static \
 && printf "#!/bin/bash\ncd /opt/sonic-annotator/\nexec ./AppRun \"\$@\"\n" >  /usr/bin/sonic-annotator \
 && find /opt/sonic-annotator -type d -exec chmod go+rx {} \; \
 && chmod go+rx \
    /usr/bin/sonic-annotator \
    /opt/sonic-annotator/AppRun \
    /opt/sonic-annotator/usr/bin/sonic-annotator

ARG BBC_VAMP_PLUGIN_VERSION=1.1

ENV VAMP_PATH /opt/sonic-annotator/usr/local/lib/vamp

RUN mkdir -p $VAMP_PATH \
 && curl -s -L -o /tmp/bbc-vamp-plugins.tar.gz \
    https://github.com/bbc/bbc-vamp-plugins/releases/download/v$BBC_VAMP_PLUGIN_VERSION/Linux.64-bit.tar.gz \
 && tar xf /tmp/bbc-vamp-plugins.tar.gz -C $VAMP_PATH \
 && rm -rf /tmp/bbc-vamp-plugins.tar.gz

COPY annotation-agent /usr/bin
COPY assets/n3 /etc/annotation-agent/n3

USER nobody

ENTRYPOINT ["/usr/bin/annotation-agent"]
