FROM postgres:16.1

RUN apt-get update && apt-get install -y cron tzdata && apt-get clean

# COPY script/backup.sh script/backup.sh
# RUN chmod +x /script/backup.sh
# RUN chmod +x /script/restore_latest.sh

RUN touch /var/log/backup.log

RUN echo "0 2 * * * /backup.sh >> /var/log/backup.log 2>&1" > /etc/cron.d/backup-cron
RUN chmod 0644 /etc/cron.d/backup-cron && crontab /etc/cron.d/backup-cron

CMD ["sh", "-c", "cron && tail -f /var/log/backup.log"]