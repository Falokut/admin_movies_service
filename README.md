# Content

+ [Configuration](#configuration)
    + [Params info](#configuration-params-info)
        + [Kafka reader config](#kafka-reader-config)
        + [time.Duration](#timeduration-yaml-supported-values)
        + [Database config](#database-config)
        + [Jaeger config](#jaeger-config)
        + [Prometheus config](#prometheus-config)
        + [Secure connection config](#secure-connection-config)
        + [Pictures config](#pictures-config)
            + [Check image config](#check-image-config)
            + [Image processing config](#image-processing-config)
+ [Metrics](#metrics)
+ [Docs](#docs)
+ [Author](#author)
+ [License](#license)
---------

# Configuration

1. [Configure movies_db](movies_db/README.md#Configuration)
2. Create .env on project root dir  
Example env:
```env
DB_PASSWORD=Password
```
3. Create a configuration file or change the config.yml file in docker\containers-configs.
If you are creating a new configuration file, specify the path to it in docker-compose volume section (your-path/config.yml:configs/)

## Configuration params info
if supported values is empty, then any type values are supported

| yml name | yml section | env name | param type| description | supported values |
|-|-|-|-|-|-|
| log_level   |      | LOG_LEVEL  |   string   |      logging level        | panic, fatal, error, warning, warn, info, debug, trace|
| host   |  listen    | HOST  |   string   |  ip address or host to listen   |  |
| port   |  listen    | PORT  |   string   |  port to listen   | The string should not contain delimiters, only the port number|
| server_mode   |  listen    | SERVER_MODE  |   string   | Server listen mode, Rest API, gRPC or both | GRPC, REST, BOTH|
| max_request_size   |  listen    | MAX_REQUEST_SIZE  |   int   | max request body size in mb, by default 4mb | only > 0|
| allowed_headers   |  listen    |  |   []string, array of strings   | list of all allowed custom headers. Need for REST API gateway, list of metadata headers, hat are passed through the gateway into the service | any strings list|
| healthcheck_port   |      | HEALTHCHECK_PORT  |   string   |     port for healthcheck       | any valid port that is not occupied by other services. The string should not contain delimiters, only the port number|
| posters_config|| |   nested yml configuration  [Pictures config](#pictures-config)   | |  |
| preview_posters_config|| |   nested yml configuration  [Pictures config](#pictures-config)   | |  |
| background_config|| |   nested yml configuration  [Pictures config](#pictures-config)   | |  |
|service_name|  prometheus    | PROMETHEUS_SERVICE_NAME | string |  service name, thats will show in prometheus  ||
|server_config|  prometheus    |   | nested yml configuration  [metrics server config](#prometheus-config) | |
|db_config|||nested yml configuration  [database config](#database-config) || configuration for database connection | |
|jaeger|||nested yml configuration  [jaeger config](#jaeger-config)|configuration for jaeger connection ||
|images_events_kafka|||nested yml configuration  [kafka reader config](#kafka-reader-config)|configuration for kafka connection||
|addr|image_storage_service|IMAGE_STORAGE_SERVICE_ADDRESS|string|ip address(or host) with port of image storage service| all valid addresses formatted like host:port or ip-address:port|
|secure_config|  image_storage_service    |   | nested yml configuration  [secure connection config](#secure-connection-config) | |
|addr|image_processing_service|IMAGE_PROCESSING_SERVICE_ADDRESS|string|ip address(or host) with port of image processing service|all valid addresses formatted like host:port or ip-address:port|
|secure_config|  image_processing_service    |   | nested yml configuration  [secure connection config](#secure-connection-config) | |


### Kafka reader config
|yml name| env name|param type| description | supported values |
|-|-|-|-|-|
|brokers||[]string, array of strings|list of all kafka brokers||
|group_id||string|id or name for consumer group||
|read_batch_timeout||time.Duration with positive duration|amount of time to wait to fetch message from kafka messages batch|[supported values](#time.Duration-yaml-supported-values)|

### time.Duration yaml supported values
A Duration value can be expressed in various formats, such as in seconds, minutes, hours, or even in nanoseconds. Here are some examples of valid Duration values:
- 5s represents a duration of 5 seconds.
- 1m30s represents a duration of 1 minute and 30 seconds.
- 2h represents a duration of 2 hours.
- 500ms represents a duration of 500 milliseconds.
- 100µs represents a duration of 100 microseconds.
- 10ns represents a duration of 10 nanoseconds.


### Database config
|yml name| env name|param type| description | supported values |
|-|-|-|-|-|
|host|DB_HOST|string|host or ip address of database| |
|port|DB_PORT|string|port of database| any valid port that is not occupied by other services. The string should not contain delimiters, only the port number|
|username|DB_USERNAME|string|username(role) in database||
|password|DB_PASSWORD|string|password for role in database||
|db_name|DB_NAME|string|database name (database instance)||
|ssl_mode|DB_SSL_MODE|string|enable or disable ssl mode for database connection|disabled or enabled|

### Jaeger config
|yml name| env name|param type| description | supported values |
|-|-|-|-|-|
|address|JAEGER_ADDRESS|string|ip address(or host) with port of jaeger service| all valid addresses formatted like host:port or ip-address:port |
|service_name|JAEGER_SERVICE_NAME|string|service name, thats will show in jaeger in traces||
|log_spans|JAEGER_LOG_SPANS|bool|whether to enable log scans in jaeger for this service or not||

### Prometheus config
|yml name| env name|param type| description | supported values |
|-|-|-|-|-|
|host|METRIC_HOST|string|ip address or host to listen for prometheus service||
|port|METRIC_PORT|string|port to listen for  of prometheus service| any valid port that is not occupied by other services. The string should not contain delimiters, only the port number|

### Secure connection config
|yml name| param type| description | supported values |
|-|-|-|-|
|dial_method|string|dial method|INSECURE,NIL_TLS_CONFIG,CLIENT_WITH_SYSTEM_CERT_POOL,SERVER|
|server_name|string|server name overriding, used when dial_method=CLIENT_WITH_SYSTEM_CERT_POOL||
|cert_name|string|certificate file name, used when dial_method=SERVER||
|key_name|string|key file name, used when dial_method=SERVER||

### Pictures config
|yml name| param type| description | supported values |
|-|-|-|-|
|validate_image|bool|validate image or not||
|check_image|nested yml configuration  [check image config](#check-image-config)|specify if validate_image is true||
|resize_image|bool|resize image or not||
|image_processing|nested yml configuration  [image processing config](#image-processing-config)|specify if resize_image is true||
|category|string|image category in storage||

#### Check image config
|yml name| param type| description | supported values |
|-|-|-|-|
|max_image_height|int32|max image height|only > 0|
|max_image_width|int32|max image width|only > 0|
|min_image_width|int32|min image width|only > 0|
|min_image_width|int32|min image width|only > 0|
|allowed_types|int32|max image height|only > 0|

#### Image processing config
|yml name| param type| description | supported values |
|-|-|-|-|
|image_height|int32|image height after resizing|only > 0|
|image_width|int32|image width  after resizing|only > 0|
|resize_method|string|resising method, by default using NearestNeighbor|Box,CatmullRom,Lanczos,Linear,MitchellNetravali|

# Metrics
The service uses Prometheus and Jaeger and supports distribution tracing

# Docs
[Swagger docs](swagger/docs/movies_service_v1.swagger.json)

# Author

- [@Falokut](https://github.com/Falokut) - Primary author of the project

# License

This project is licensed under the terms of the [MIT License](https://opensource.org/licenses/MIT).

---