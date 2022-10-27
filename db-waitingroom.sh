#!/bin/sh

set -e
COMMAND=$@

echo 'Waiting for database to be available...'
maxTries=10
while [ "$maxTries" -gt 0 ] && ! mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DB" -e 'SHOW TABLES' > /dev/null 2>&1; do
    maxTries=$(($maxTries - 1))
    echo 'Retrying to connect to MySQL Database...'
    sleep 3
done
echo
if [ "$maxTries" -le 0 ]; then
    echo >&2 'error: unable to contact mysql after 10 tries'
    exit 1
fi

exec $COMMAND