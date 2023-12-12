# logs-microservices

* All microservices are written in Go language
* Logs are generate from **backend/log_generator**
* Logs are read from file by **backend/microservices/log_collector** and then put into **Kafka** Topic
* Logs are read from Kafka by **backend/microservices/log_indexer**

# Kafka

* Kafka config is defined under **externals/kafka**
* In the docker-compose file we also build the **kowl** image to see the logs that are on the Kafka topic (Note that it might take about 30 seconds to 1 minute to be visible due to networking)
