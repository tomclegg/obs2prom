# obs2prom
Export [prometheus](https://github.com/prometheus/prometheus) metrics for your [open broadcaster](https://github.com/openbroadcaster/obplayer)'s emergency alert system... so when your alerter can alert you when it can't get alerts.

Quick start:
* Install Go tools, set $GOPATH
* `go get github.com/tomclegg/obs2prom`
* `install $GOPATH/bin/obs2prom /usr/local/bin/`
* `cp $GOPATH/src/github.com/tomclegg/obs2prom/obs2prom.service /lib/systemd/system/`
* Edit `/lib/systemd/system/obs2prom.service` -- e.g., add command line arguments to the ExecStart line like `-alerts.url http://username:password@obhost:23233/alerts/list`
* `systemctl enable obs2prom`
* `systemctl start obs2prom`

Add a prometheus target like `localhost:9911`.

Add graphs/alerts for `round(time()-obs_alert_heartbeat_time_seconds)`.

Run `obs2prom -help` for other options.
