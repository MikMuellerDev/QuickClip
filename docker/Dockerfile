FROM amd64/debian:stable-slim

ADD ./QuickClip /app/quickclip
ADD ./run.sh /app/quickclip/bin/run.sh

WORKDIR /app/quickclip/

RUN export DEBIAN_FRONTEND=noninteractive && apt update && apt install -yq \
  tzdata \
  curl \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s CMD curl --fail http://localhost:80/health || exit 1
EXPOSE 80
ENTRYPOINT ["/app/quickclip/bin/run.sh"]