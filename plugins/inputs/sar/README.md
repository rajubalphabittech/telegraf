# Sar Input Plugin

The sar plugin gathers metrics from **sar**(from sysstat utilities family) output files.
It runs a **sar** command(`sar -f <sa_file>`) to read the output of file.
The utility binary is configurable from the config.
The plugin looks up for files in `inputs_dir_path` directory that can be configured.
After processing the file it moves into another directory that is also configured from config file(`outputs_dir_path`)


### Configuration:

```toml
## Path to the sar command. 
#
sar_path = "usr/bin/sar"
#
#
## Input files directory path
#
inputs_dir_path = "./sar_inputs"
#
#
## Output file path
outputs_dir_path = "./sar_processed"
```

### Measurements & Fields:

Measurements depend on the flags that are set to the utility.
By default, the output files will have 

- sar 
    - system (float, percent)
    - iowait (float, percent)
    - steal (float, percent)
    - idle (float, percent)
    - CPU (string, unit)
    - user (float, percent)
    - nice (float, percent)
- average 
    - system (float, percent)
    - iowait (float, percent)
    - steal (float, percent)
    - idle (float, percent)
    - CPU (string, unit)
    - user (float, percent)
    - nice (float, percent)

### Tags:

- All measurements have the following tags:
  - linux_kernel
  - host
  - date
  - arch
  - num_cpu

### Example Output:

```
[max@thinkmaxthink telegraf]$ telegraf --config=telegraf.conf --test
* Plugin: inputs.sar, Collection 1
> sar,arch=_x86_64_,num_cpu=4,linux_kernel=Linux4.8.15-200.fc24.x86_64,host=(thinkmaxthink),date=06/04/2017 system=1.75,iowait=0,steal=0,idle=93.98,CPU="all",user=2.76,nice=1.5 -6826973833871345152
> sar,linux_kernel=Linux4.8.15-200.fc24.x86_64,host=(thinkmaxthink),date=06/04/2017,arch=_x86_64_,num_cpu=4 CPU="all",user=3.02,nice=2.27,system=2.02,iowait=0,steal=0,idle=92.7 -6826973832871345152
> sar,linux_kernel=Linux4.8.15-200.fc24.x86_64,host=(thinkmaxthink),date=06/04/2017,arch=_x86_64_,num_cpu=4 idle=93.97,CPU="all",user=2.76,nice=1.51,system=1.76,iowait=0,steal=0 -6826973831871345152
> sar,date=06/04/2017,arch=_x86_64_,num_cpu=4,linux_kernel=Linux4.8.15-200.fc24.x86_64,host=(thinkmaxthink) CPU="all",user=3.28,nice=1.89,system=2.43,iowait=0.11,steal=0,idle=92.28 -6826973820871345152
> sar,host=(thinkmaxthink),date=06/04/2017,arch=_x86_64_,num_cpu=4,linux_kernel=Linux4.8.15-200.fc24.x86_64 nice=1.51,system=1.76,iowait=0,steal=0,idle=93.95,CPU="all",user=2.77 -6826973819871345152
> sar,num_cpu=4,linux_kernel=Linux4.8.15-200.fc24.x86_64,host=(thinkmaxthink),date=06/04/2017,arch=_x86_64_ user=2.76,nice=2.01,system=2.51,iowait=0,steal=0,idle=92.73,CPU="all" -6826973818871345152
> sar,num_cpu=4,linux_kernel=Linux4.8.15-200.fc24.x86_64,host=(thinkmaxthink),date=06/04/2017,arch=_x86_64_ CPU="all",user=2.26,nice=1.5,system=1.75,iowait=0,steal=0,idle=94.49 -6826973817871345152
> sar,linux_kernel=Linux4.8.15-200.fc24.x86_64,host=(thinkmaxthink),date=06/04/2017,arch=_x86_64_,num_cpu=4 iowait=0.25,steal=0,idle=93.72,CPU="all",user=2.51,nice=1.26,system=2.26 -6826973816871345152
> sar,linux_kernel=Linux4.8.15-200.fc24.x86_64,host=(thinkmaxthink),date=06/04/2017,arch=_x86_64_,num_cpu=4 nice=2.52,system=1.76,iowait=0,steal=0,idle=92.7,CPU="all",user=3.02 -6826973815871345152

```

### Installation
1. Unpack the zip file, the contents must be `sar` directory.
2. Copy the directory into `telegraf/plugin/inputs/` directory.
3. Add the plugin into telegaf plugins registy:
    1. Open the file: `telegraf/plugins/inputs/all/all.go`
    2. Add the following into imports: `_ "github.com/influxdata/telegraf/plugins/inputs/sar"`
4. Rebuild the **Telegraf**
5. Configure the plugin and run.

#### Example Telegraf config:
```
# Telegraf Configuration
#
# Telegraf is entirely plugin driven. All metrics are gathered from the
# declared inputs, and sent to the declared outputs.
#
# Plugins must be declared in here to be active.
# To deactivate a plugin, comment out the name and any variables.
#
# Use 'telegraf -config telegraf.conf -test' to see what metrics a config
# file would generate.
#
# Environment variables can be used anywhere in this config file, simply prepend
# them with $. For strings the variable must be within quotes (ie, "$STR_VAR"),
# for numbers and booleans they should be plain (ie, $INT_VAR, $BOOL_VAR)


# Global tags can be specified here in key="value" format.
[global_tags]
  # dc = "us-east-1" # will tag all metrics with dc=us-east-1
  # rack = "1a"
  ## Environment variables can be used as tags, and throughout the config file
  # user = "$USER"


# Configuration for telegraf agent
[agent]
  ## Default data collection interval for all inputs
  interval = "10s"
  ## Rounds collection interval to 'interval'
  ## ie, if interval="10s" then always collect on :00, :10, :20, etc.
  round_interval = true

  ## Telegraf will send metrics to outputs in batches of at most
  ## metric_batch_size metrics.
  ## This controls the size of writes that Telegraf sends to output plugins.
  metric_batch_size = 1000

  ## For failed writes, telegraf will cache metric_buffer_limit metrics for each
  ## output, and will flush this buffer on a successful write. Oldest metrics
  ## are dropped first when this buffer fills.
  ## This buffer only fills when writes fail to output plugin(s).
  metric_buffer_limit = 10000

  ## Collection jitter is used to jitter the collection by a random amount.
  ## Each plugin will sleep for a random time within jitter before collecting.
  ## This can be used to avoid many plugins querying things like sysfs at the
  ## same time, which can have a measurable effect on the system.
  collection_jitter = "0s"

  ## Default flushing interval for all outputs. You shouldn't set this below
  ## interval. Maximum flush_interval will be flush_interval + flush_jitter
  flush_interval = "10s"
  ## Jitter the flush interval by a random amount. This is primarily to avoid
  ## large write spikes for users running a large number of telegraf instances.
  ## ie, a jitter of 5s and interval 10s means flushes will happen every 10-15s
  flush_jitter = "0s"

  ## By default or when set to "0s", precision will be set to the same
  ## timestamp order as the collection interval, with the maximum being 1s.
  ##   ie, when interval = "10s", precision will be "1s"
  ##       when interval = "250ms", precision will be "1ms"
  ## Precision will NOT be used for service inputs. It is up to each individual
  ## service input to set the timestamp at the appropriate precision.
  ## Valid time units are "ns", "us" (or "µs"), "ms", "s".
  precision = ""

  ## Logging configuration:
  ## Run telegraf with debug log messages.
  debug = false
  ## Run telegraf in quiet mode (error log messages only).
  quiet = false
  ## Specify the log file name. The empty string means to log to stderr.
  logfile = ""

  ## Override default hostname, if empty use os.Hostname()
  hostname = ""
  ## If set to true, do no set the "host" tag in the telegraf agent.
  omit_hostname = false


###############################################################################
#                            OUTPUT PLUGINS                                   #
###############################################################################

# Configuration for influxdb server to send metrics to
[[outputs.influxdb]]
  ## The HTTP or UDP URL for your InfluxDB instance.  Each item should be
  ## of the form:
  ##   scheme "://" host [ ":" port]
  ##
  ## Multiple urls can be specified as part of the same cluster,
  ## this means that only ONE of the urls will be written to each interval.
  # urls = ["udp://localhost:8089"] # UDP endpoint example
  urls = ["http://localhost:8086"] # required
  ## The target database for metrics (telegraf will create it if not exists).
  database = "telegraf" # required

  ## Name of existing retention policy to write to.  Empty string writes to
  ## the default retention policy.
  retention_policy = ""
  ## Write consistency (clusters only), can be: "any", "one", "quorum", "all"
  write_consistency = "any"

  ## Write timeout (for the InfluxDB client), formatted as a string.
  ## If not provided, will default to 5s. 0s means no timeout (not recommended).
  timeout = "5s"
  # username = "telegraf"
  # password = "metricsmetricsmetricsmetrics"
  ## Set the user agent for HTTP POSTs (can be useful for log differentiation)
  # user_agent = "telegraf"
  ## Set UDP payload size, defaults to InfluxDB UDP Client default (512 bytes)
  # udp_payload = 512

  ## Optional SSL Config
  # ssl_ca = "/etc/telegraf/ca.pem"
  # ssl_cert = "/etc/telegraf/cert.pem"
  # ssl_key = "/etc/telegraf/key.pem"
  ## Use SSL but skip chain & host verification
  # insecure_skip_verify = false


###############################################################################
#                            INPUT PLUGINS                                    #
###############################################################################

# # Sar utility output collector
[[inputs.sar]]
sar_path = "/usr/bin/sar"
inputs_dir_path = "/tmp/sar/"
outputs_dir_path = "/tmp/sar_processed/"
```
