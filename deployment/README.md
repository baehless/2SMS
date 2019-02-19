# Deployment of 2SMS
The scripts in this folder can be used to initially deploy the different components of 2SMS. All the scripts can be adapted (or
even be used directly) to perform an update of an existing deployment.
There are some points that should be noted:
- The current system fully supports only a single Scraper instance (multiple instances can be deployed but are not interconnected)
- In the default setup Manager and Scraper are installed in the same local network behind a NAT, other setups are possible but
    require adapting the install scripts (application's parameters)
- If a deployment from scratch is required, then some files in the downloaded configuration archives must be updated (see below for more details)
- Deployment of other external components (e.g. alertmanager, blackbox_exporter, node_exporter, ...) have to be done manually (see below)
- At this stage understanding of Prometheus other than its query language is still necessary as different operations have to
    be performed manually. The documentation can be found at https://prometheus.io/docs/introduction/overview/

## Requirements
- Every instance needs to be reachable via a public IPv4 address
- Every installation scripts will require sudo right to run systemctl commands
- SCION must be installed and a connection to SCIONLab is required
    - $SC/gen/ia must exist
    - /run/shm/sciond/default.sock must exist
    - the AS must be set up as endhost

## Procedure to deploy from scratch
    1) Run `install_manager.sh` (change IP accordingly to the setting) to install and run the Manager. Note that this
        will create a new CA with corresponding root certificate and therefore all old issued certificates (e.g. to Endpoints)
        will not be trusted anymore. If you already have a 2SMS installation and wish to keep the same root certificate then
        you need to copy the old ca directory in the new installation path (default to ~/2SMS/deployment/manager)
    2) Replace the following files in scraper_configuration/ca_certs and endpoint_configuration/ca_certs with the ones
        that are generated in manager/ca:
            - bootstrap.json
            - ca.crt
       Replace the AS certificate in those two directories with the (latest) one in the manager/ca_certs directory
    3) Compress (.tar.gz) the modified directory and load them on the website
    4) Run `curl localhost:10002/manager/signing/enable` to enable creation of new certificates for new Scraper or Endpoint instances
        that will be created next (otherwise the component will block and periodically try to get a certificate)
    5) Run `install_scraper.sh` (change MANAGER_IP and other parameters accordingly to the setting) to install and run the Scraper.
    6) Run `install_endpoint.sh` (change MANAGER_IP accordingly to your setting) in another machine to install and run an Endpoint instance.
    7) When you are finished run `curl localhost:10002/manager/signing/block` to block creation of new certificates at the manager

## Procedure to deploy on the existing SCIONLab infrastructure
TODO

# Deployment of external components

## Procedure to install alertmanager
TODO

## Procedure to install blackbox exporter
TODO

## Procedure to install node exporter
TODO

## Procedure to install InfluxDB
TODO


# FAQ
## How do I add a new alert?
-> To add a new alert it's necessary to change `scraper/prometheus/alert_rules.yml`.
    A simple alert has this format:
    - alert: <alert_name>               # A unique name for the alert
        expr: <triggering_expression>   # The expression (PromQL syntax) that will trigger the alert when evaluated to true
        for: <duration>                 # For how long the expression has to evaluate to true before the alert is really fired
    More information about syntax and options can be found at https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/
    Finally the Prometheus server must be reloaded for the changes to be applied. In a basic installation his can be done with the following command:
    `curl -X POST localhost:9090/-/reload`. In a custom installation the URL may be different (e.g. see reload_prometheus.sh in PVM).

## How is a target configured in the Prometheus' server configuration file?
-> Every target is configured as a job with this format: TODO
    - job_name: <IA> <IP <type>             # A unique name identifying the target (e.g. 17-ffaa:1:c5 127.0.0.1 bs)
      metrics_path: /<IA>/<local_path>      # (e.g. /17-ffaa:1:c5/bs)
      static_configs:                       #
      - targets:                            #
        - <IP>:<Port> (e.g. 127.0.0.1:9199) #
        labels:                             #
          AS: <AS> (e.g. ffaa:1:c5)         #
          ISD: "17"                         #
          service: bs                       #
  proxy_url: http://127.0.0.1:9901          #
  More information about options and syntax can be found at https://prometheus.io/docs/prometheus/latest/configuration/configuration/

## I started a new Endpoint but it's not being monitored, what went wrong?
-> There are multiple possible reasons for this behaviour:
    - Signing of new certificates at the manager is blocked and the Endpoint blocks until it is able to get one
    - TODO