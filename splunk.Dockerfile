FROM splunk/splunk

COPY ./interlock_mon.xml /opt/splunk/etc/apps/search/local/data/ui/views/interlock_mon.xml