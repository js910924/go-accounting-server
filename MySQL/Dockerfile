FROM mysql

# Add a database
ENV MYSQL_DATABASE account

COPY ./sql-scripts/ /docker-entrypoint-initdb.d/
