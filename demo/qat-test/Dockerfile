FROM debian:sid as builder 
ENV http_proxy http://proxy.ir.intel.com:911/
ENV https_proxy http://proxy.ir.intel.com:911/
ENV no_proxy http://proxy.ir.intel.com:911

ENV QAT_DRIVER_RELEASE="qat1.7.l.4.4.0-00023"
ENV QAT_ENGINE_VERSION="v0.5.39"
ENV PERL5LIB="/OSL/openssl"
ENV NGINX_INSTALL_DIR="/NGINX/asynch_mode_nginx"
ENV NGINX_SRC_DIR="/NGINX/asynch_mode_nginx"
ENV OPENSSL_LIB="/usr/local/ssl"
ENV OPENSSL_SRC_DIR="OSL/openssl"
ENV OPENSSL_ENGINES="/usr/local/ssl/lib/engines-1.1"
ENV QZ_ROOT="/QATZ/QATzip"
ENV QATZIP_SRC_DIR="/QATZ/QATzip"
ENV ICP_ROOT="/QAT"
ENV LD_LIBRARY_PATH="/OSL/openssl/lib:/QAT/build:/QATZ/QATzip/src/"  
       

RUN apt-get update && \
    apt-get install -y git build-essential wget libssl-dev openssl libudev-dev pkg-config autoconf autogen libtool libpcre3 libpcre3-dev libxml2-dev libxslt1-dev python-dev zlib1g zlib1g-dev procps vim
RUN mkdir /QAT && mkdir /QATE && mkdir /QATZ && mkdir /OSL && mkdir /NGINX 
WORKDIR /QAT
RUN wget https://01.org/sites/default/files/downloads//$QAT_DRIVER_RELEASE.tar.gz && \
    tar zxf $QAT_DRIVER_RELEASE.tar.gz
  
RUN sed -i -e 's/cmn_ko$//' -e 's/lac_kernel$//' quickassist/Makefile && \
    KERNEL_SOURCE_ROOT=/tmp ./configure && \
    make quickassist-all adf-ctl-all && \
    install -m 755 build/libqat_s.so /usr/lib/ && \
    install -m 755 build/libusdm_drv_s.so /usr/lib/ && \
    install -m 755 build/adf_ctl /usr/bin/ 

WORKDIR /OSL
RUN git clone -b OpenSSL_1_1_1a https://github.com/openssl/openssl.git
WORKDIR /OSL/openssl
RUN ./config --prefix=/usr/local/ssl -Wl,-rpath,\${LIBRPATH}
RUN make && make install
WORKDIR /usr/local/ssl/lib
#RUN mkdir engines-1.1
WORKDIR /QATE
RUN git clone -b v0.5.40 https://github.com/intel/QAT_Engine
WORKDIR /QATE/QAT_Engine
RUN ./autogen.sh && \
    ./configure \
    --with-qat_dir=/QAT \
    --with-openssl_dir=/OSL/openssl \
    --with-openssl_install_dir=/usr/local/ssl \
    --enable-upstream_driver \
    --enable-usdm \
    --with-qat_install_dir=/usr/lib && \
    make && make install

WORKDIR /QATZ
RUN git clone -b v0.2.7 https://github.com/intel/QATzip
WORKDIR /QATZ/QATzip
RUN ./setenv.sh && \
    ./configure --with-ICP_ROOT=/QAT && \
    make clean && make install

WORKDIR /NGINX
RUN git clone -b v0.3.7 https://github.com/intel/asynch_mode_nginx
WORKDIR /NGINX/asynch_mode_nginx
RUN sed -e '/--with-http_xslt_module/ s/^#*/#/' -i test/nginx-config.test && \
    ./configure \
     --prefix=$NGINX_INSTALL_DIR \
     --with-http_ssl_module \
     --add-dynamic-module=modules/nginx_qat_module/ 
RUN  sed -e '/\tcp conf/ s/^#*/#/' -i objs/Makefile 
RUN  make && make install

#FROM debian:sid-slim
#ENV http_proxy http://proxy.ir.intel.com:911/
#ENV https_proxy http://proxy.ir.intel.com:911/
#ENV no_proxy http://proxy.ir.intel.com:911

#RUN apt-get update && apt-get install -y openssl

#COPY --from=builder /usr/lib/libqat_s.so /usr/lib/
#COPY --from=builder /usr/lib/libusdm_drv_s.so /usr/lib/
#COPY --from=builder /usr/bin/adf_ctl /usr/bin
#COPY --from=builder /usr/lib/x86_64-linux-gnu/engines-1.1/qat.so /usr/lib/x86_64-linux-gnu/engines-1.1/qat.so

