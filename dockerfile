FROM scratch
ADD application /
ADD config.json /
ADD index.html /
WORKDIR /
CMD ["/application"]
EXPOSE 5000
