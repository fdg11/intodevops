FROM fdg11/nodejs-gulp:latest
MAINTAINER Alex P "alexfile123@gmail.com"

# Install general Dependencies
#RUN apt-get update && apt-get install -yqq \
#      curl \
#      ca-certificates \
#      build-essential \
#      && apt-get autoremove -y \
#      && apt-get autoclean

# Install node.js
#RUN curl -sL https://deb.nodesource.com/setup_8.x | bash -
#RUN apt-get update && apt-get install -yqq \
#     nodejs \
#      && apt-get autoremove -y \
#      && apt-get autoclean

## NPM, GULP 
# Global install gulp
<<<<<<< 715e64a66d738b67d147cfd2087f234b031edbef
#RUN npm install -g gulp
=======
#RUN npm install gulp-autoprefixer
RUN npm install -g gulp
>>>>>>> update

# Define working directory.
#WORKDIR /workspace

ADD ./install.sh /script/install.sh

ENTRYPOINT ["/script/install.sh"]

EXPOSE 3000

