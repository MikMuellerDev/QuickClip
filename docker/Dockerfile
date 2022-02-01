FROM amd64/debian:stable-slim

ADD ./QuickClip /app/QuickClip
ADD ./run.sh /app/QuickClip/bin/run.sh

WORKDIR /app/QuickClip/

RUN export DEBIAN_FRONTEND=noninteractive && apt update && apt install -yq \
  tzdata \
  curl \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s CMD curl --fail http://localhost:8081/health || exit 1
EXPOSE 8081
ENTRYPOINT ["/app/QuickClip/bin/run.sh"]