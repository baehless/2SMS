# TODO: move to readme
# 2SMS job format:
#- job_name: <IA> <IP <type> (e.g. 17-ffaa:1:c5 127.0.0.1 bs)
#  metrics_path: /<IA>/<local_path> (e.g. /17-ffaa:1:c5/bs)
#  static_configs:
#  - targets:
#    - <IP>:<Port> (e.g. 127.0.0.1:9199)
#    labels:
#      AS: <AS> (e.g. ffaa:1:c5)
#      ISD: "17"
#      service: bs
#  proxy_url: http://127.0.0.1:9901