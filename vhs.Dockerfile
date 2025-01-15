FROM ghcr.io/charmbracelet/vhs:v0.9.0 as vhs

RUN export DEBIAN_FRONTEND=noninteractive && \
  apt-get update && \
  apt-get install -y --no-install-recommends \
  ca-certificates \
  curl \
  poppler-utils && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*

# Install go
COPY --from=golang:bookworm /usr/local/go /usr/local/go
ENV PATH="/usr/local/go/bin:${PATH}"

# Install task
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin && \
  task --version

COPY . /usr/src/app
WORKDIR /usr/src/app

RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg/mod \
  task build && \
  cp papercrypt /usr/local/bin/papercrypt && \
  papercrypt version

WORKDIR /vhs
