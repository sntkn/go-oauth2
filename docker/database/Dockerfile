FROM postgres:latest

ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=admin

# Add initialization scripts to the Docker entrypoint directory
COPY ./init /docker-entrypoint-initdb.d

RUN sleep 60

# Expose the PostgreSQL port
EXPOSE 5432