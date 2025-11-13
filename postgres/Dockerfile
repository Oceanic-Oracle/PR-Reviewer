FROM postgres:15-alpine

ENV POSTGRES_DB=postgres
ENV POSTGRES_USER=User

COPY migrations/ /docker-entrypoint-initdb.d/

EXPOSE 5432