# Choose the base image for the container
FROM ubuntu:latest

# Copy the binary file and the script from the host system into the container
COPY ./order /app/order
COPY ./limit-bandwidth.sh /app/limit-bandwidth.sh

# Change permissions for the binary and the script
RUN chmod +x /app/order
RUN chmod +x /app/limit-bandwidth.sh

# Run the script and the program when the container starts
CMD /app/limit-bandwidth.sh && /app/order