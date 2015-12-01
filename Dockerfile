FROM scratch

# Copy executable and config
COPY wedirect /

# If does not exist create from config.json.example
COPY config.json /

# Container settings
EXPOSE 8080
ENTRYPOINT ["/wedirect"]
