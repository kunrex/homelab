FROM ubuntu:24.04

# Set common environment variables
ENV DEBIAN_FRONTEND=noninteractive

ENV LANG=en_US.UTF-8
ENV LANGUAGE=en_US.UTF-8

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install core
RUN apt-get update && apt-get -y upgrade \
    && apt-get install -y openssh-server \
    # Utils
    && apt-get install -y mc btop iotop ncdu tar zip nano vim bash sudo sed wget ca-certificates curl unzip gnupg build-essential git ninja-build gettext cmake fd-find ripgrep gh age byobu software-properties-common bat tmux \
    # Python
    && apt-get install -y python3 python3-pip python3-venv \
    # Go
    && apt-get install -y golang-go \
    # .NET and Mono
    && apt-get install -y dotnet-sdk-10.0 aspnetcore-runtime-10.0 mono-complete \
    # Net utils
    && apt-get install -y iputils-ping traceroute telnet dnsutils iperf nmap

ENV PATH="/usr/local/go:$PATH"

# Install docker client
RUN install -m 0755 -d /etc/apt/keyrings \
    && curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg \
    && chmod a+r /etc/apt/keyrings/docker.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo $VERSION_CODENAME) stable" \
        > /etc/apt/sources.list.d/docker.list \
    && apt-get update \
    && apt-get install -y \
        redis-tools \
        postgresql-client \
        docker-ce \
        docker-ce-cli \
        docker-compose-plugin \
        docker-buildx-plugin

# Install fastfetch
RUN add-apt-repository ppa:zhangsongcui3371/fastfetch \
    && apt-get update && apt-get install -y fastfetch

# Install Node.js
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get update && apt-get install -y nodejs \
    && npm install -g npm@latest

# Install nvim
RUN curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim-linux-x86_64.tar.gz \
    && rm -rf /opt/nvim-linux-x86_64 \
    && tar -C /opt -xzf nvim-linux-x86_64.tar.gz \
    && rm nvim-linux-x86_64.tar.gz

ENV PATH="/opt/nvim-linux-x86_64/bin:$PATH"

# Config SSH
# Uncomment root login and disallow it, ssh access is given through Tailscale
RUN sed -ri "s|^#PermitRootLogin|PermitRootLogin|" /etc/ssh/sshd_config \
    && sed -i "s|PermitRootLogin without-password|PermitRootLogin yes|" /etc/ssh/sshd_config \
    && sed -i "s|PermitRootLogin prohibit-password|PermitRootLogin yes|" /etc/ssh/sshd_config \
    && sed -ri "s|^#?PermitRootLogin\s+.*|PermitRootLogin no|" /etc/ssh/sshd_config \
    && sed -ri "s|UsePAM yes|#UsePAM yes|g" /etc/ssh/sshd_config

# Cleaning
RUN apt-get clean autoclean -y \
    && apt-get autoremove -y \
    && rm -rf /var/lib/{apt,dpkg,cache,log}/ \
    && rm -rf /var/lib/apt/lists/*.lz4 \
    && rm -rf /var/log/* \
    && rm -rf /tmp/* \
    && rm -rf /var/tmp/* \
    #    && rm -rf /usr/share/doc/ \
    #    && rm -rf /usr/share/man/ \
    && rm -rf $HOME/.cache

# Install Yazi
RUN wget https://github.com/sxyazi/yazi/releases/latest/download/yazi-x86_64-unknown-linux-gnu.zip \
    && unzip yazi-x86_64-unknown-linux-gnu.zip -d yazi-temp \
    && mv yazi-temp/*/yazi yazi-temp/*/ya /usr/local/bin/ \
    && rm -rf yazi-temp yazi-x86_64-unknown-linux-gnu.zip

#Install Zoxide
RUN curl -sSfL https://raw.githubusercontent.com/ajeetdsouza/zoxide/main/install.sh | sh -s - --bin-dir /usr/local/bin/

# Install fzf
RUN curl -L https://github.com/junegunn/fzf/releases/download/0.48.0/fzf-0.48.0-linux_amd64.tar.gz -o /tmp/fzf.tar.gz \
    && tar -xzf /tmp/fzf.tar.gz -C /usr/local/bin \ 
    && rm /tmp/fzf.tar.gz

# Install oh-my-posh
RUN curl -s https://ohmyposh.dev/install.sh | bash -s -- -d /usr/local/bin

# Install ctop
RUN wget https://github.com/bcicen/ctop/releases/download/v0.7.7/ctop-0.7.7-linux-amd64 -O /usr/local/bin/ctop \
    && chmod +x /usr/local/bin/ctop

ARG WORKSPACE_UID
ARG WORKSPACE_GID
ARG WORKSPACE_USER
ARG WORKSPACE_PASSWORD
ENV WORKSPACE_USER=$WORKSPACE_USER

ARG DOCKER_GID

RUN groupadd -g $WORKSPACE_GID $WORKSPACE_USER \
    && useradd -m -s /bin/bash -u $WORKSPACE_UID -g $WORKSPACE_GID $WORKSPACE_USER \
    && echo "${WORKSPACE_USER}:${WORKSPACE_PASSWORD}" | chpasswd 

# Optional, give sudo and docker permissions
RUN groupmod -g $DOCKER_GID docker \
    && usermod -aG sudo $WORKSPACE_USER \
    && usermod -aG $DOCKER_GID $WORKSPACE_USER 

COPY ./scripts /home/$WORKSPACE_USER
RUN chown -R $WORKSPACE_USER:$WORKSPACE_USER /home/$WORKSPACE_USER

USER $WORKSPACE_USER
WORKDIR /home/$WORKSPACE_USER

# Prepare SSH configuration
RUN mkdir -p /home/$WORKSPACE_USER/.ssh \
    && touch /home/$WORKSPACE_USER/.ssh/known_hosts

# npm setup
RUN mkdir -p /home/$WORKSPACE_USER/.npm-global \
    && chown -R $WORKSPACE_USER:$WORKSPACE_USER /home/$WORKSPACE_USER/.npm-global \
    && npm config set prefix "/home/$WORKSPACE_USER/.npm-global"

ENV PATH="/home/$WORKSPACE_USER/.npm-global/bin:$PATH"

# Install pnpm
RUN npm install -g pnpm

ENV PATH="/home/$WORKSPACE_USER/.venv/bin:$PATH"
ENV GOPATH="/home/$WORKSPACE_USER/go"
ENV PATH="$GOPATH/bin:$PATH"

# Install OpenCode
RUN npm install -g opencode-ai

# Add PATH to path
RUN echo "export PATH=\$PATH:$PATH:" >> /home/$WORKSPACE_USER/.path

USER root

COPY entry/entrypoint.sh /
RUN chmod +x /entrypoint.sh

EXPOSE 22/tcp
ENTRYPOINT ["/entrypoint.sh"]

CMD ["/usr/sbin/sshd", "-D"]
