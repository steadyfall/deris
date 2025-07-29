FROM golang:1.24.5-bookworm

RUN apt-get update && apt-get install -y \
	build-essential \
	cmake \
	git \
	libjson-c-dev \
	libwebsockets-dev \
	ffmpeg \
	tmux \
	chromium \
	sudo \
	telnet \
	vim \
	&& rm -rf /var/lib/apt/lists/*

# ttyd not installed directly on debian-based distros
# ref: https://github.com/charmbracelet/vhs/issues/527
RUN cd /root \
	&& git clone https://github.com/tsl0922/ttyd.git \
	&& cd ttyd && mkdir build && cd build \
	&& cmake .. \
	&& make && sudo make install

# chromium cannot be run as root
# ref: https://issues.chromium.org/issues/40480798
RUN useradd -m nonroot \
	&& echo "PS1='\[\e[38;5;49m\]\$\[\e[0m\] '" >> /home/nonroot/.bashrc \
	&& chown -R nonroot:nonroot /home/nonroot/

RUN go install github.com/charmbracelet/vhs@latest

USER nonroot
WORKDIR /home/nonroot