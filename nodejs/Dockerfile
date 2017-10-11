FROM ubuntu:latest
MAINTAINER Alex P "alexfile123@gmail.com"

# Install general Dependencies
RUN apt-get update && apt-get install -yqq \
      curl \
      ca-certificates \
      build-essential \
      && apt-get autoremove -y \
      && apt-get autoclean

# Install node.js
RUN curl -sL https://deb.nodesource.com/setup_8.x | bash -
RUN apt-get update && apt-get install -yqq \
      nodejs \
      && apt-get autoremove -y \
      && apt-get autoclean

## NPM, GULP 
# Global install gulp
RUN npm install -g gulp

# Define working directory.
WORKDIR /workspace

EXPOSE 3000
