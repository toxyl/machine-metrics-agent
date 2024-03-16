# Machine Metrics Agent (MMA)
This is a small application that collects metrics from a host and writes them to an Influx DB.

## Metrics Collected
| Group | Metric | Unit | Desc |
| --- | --- | --- | --- | 
| cpu | cores | int | Cores available |
| cpu | used_percent | % (0..1) | Used CPU (measured over all cores) |
| mem | total | Bytes | Total memory |
| mem | avail | Bytes | Available memory |
| mem | avail_percent | % (0..1) | % available memory  |
| mem | used | Bytes | Used memory |
| mem | used_percent | % (0..1) | % used memory |
| disk | total | Bytes | Total disk space |
| disk | avail | Bytes | Available disk space |
| disk | avail_percent | % (0..1) | % available disk space |
| disk | used | Bytes | Used disk space |
| disk | used_percent | % (0..1) | % used disk space |
| load | avg1 | Bytes | 1m load average |
| load | avg1_percent | % (0..1) | 1m load average (measured over all cores) |
| load | avg5 | Bytes | 5m load average |
| load | avg5_percent | % (0..1) | 5m load average (measured over all cores) |
| load | avg15 | Bytes | 15m load average |
| load | avg15_percent | % (0..1) | 15m load average (measured over all cores) |
| net | in | Bytes | Total received |
| net | out | Bytes | Total sent |
| uptime | seconds | seconds | Amount of seconds since last boot |