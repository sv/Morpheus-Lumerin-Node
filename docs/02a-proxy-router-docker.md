
# Proxy-Router Docker Setup: 

**This document describes setting up the proxy-router component of the Morpheus AI Network in a Docker container ONLY and accessing it via the Swagger API Interface...no GUI or Wallet components are included in this setup**

## Overview: 
* The Proxy-Router is a critical component of the Morpheus AI Network. It is responsible for monitoring the blockchain for events and managing the AI models and providers that are available to the network. 
* The Proxy-Router is a standalone application that can be run on any server that has access to the blockchain and the AI models.
* This document walks through using Docker to build and run the proxy-router image on your server.

## Pre-Requisites: 
* Your AI model has been configured, started and made available to the proxy-router server via a private endpoint (IP:PORT or DNS:PORT) eg: `http://mycoolaimodel.domain.com:8080`
* You have an existing funded wallet with MOR and ETH and also have the `private key` for the wallet (this will be needed for the .env file configuration)
* Your proxy-router must have a **publicly accessible endpoint** for the provider (ip:port or fqdn:port no protocol) eg: `mycoolmornode.domain.com:3333` - this will be used when creating the provider on the blockchain
* Docker and Git are installed and current on your server
* The three key configuration files are located in the same directory as the proxy-router code: 
    * .env
    * models-config.json
    * rating-config.json

## Installation & Configuration Steps:
1. Clone the source code from the repository: 
    * `git clone -b main https://github.com/Lumerin-protocol/Morpheus-Lumerin-Node.git`  
1. Change to the proxy-router directory:
    * `cd Morpheus-Lumerin-Node\proxy-router`
1. Configure the 3 key files for your environment:
    1. Environment file configuration 
        1. Copy the example file to .env: 
            * Linux/Mac: `cp .env.example .env`  
            * Windows: `copy .env.example.win .env`
        1. Edit values within the .env as desired
            * Add your private key to`WALLET_PRIVATE_KEY=`
            * Modify the following values to ensure that those files remain "outside of the running container" for persistence and are mounted by the docker-compose.yml file's `volume` directive
                * `MODELS_CONFIG_PATH=/app/data/models-config.json` 
                * `RATING_CONFIG_PATH=/app/data/rating-config.json`
                * `PROXY_STORAGE_PATH=/app/data/`
        1. Choose the **blockchain** you'd like to work on...**Arbitrum MAINNET is the default** 
            * To operate on the Sepolia Arbitrum TESTNET,  
            * Uncomment the `TESTNET VALUES` and comment the `MAINNET VALUES` lines & save the file
    1. Models Configuration file
        1. Copy the example file to models-config.json: 
            * Linux/Mac: `cp models-config.json.example models-config.json`  
            * Windows: `copy models-config.json.example models-config.json` 
        1. Edit the models-config.json file to include the models you'd like to use. 
        1. Details here: [models-config.json.md](models-config.json.md)
        1. Once your provider is up and running, deploy a new model and model bid via the API interface (you will need to update the `modelId` for the configuration)
    1. Rating Configuration file
        1. Copy the example file to rating-config.json: 
            * Linux/Mac: `cp rating-config.json.example rating-config.json`  
            * Windows: `copy rating-config.json.example rating-config.json` 
        1. Edit the rating-config.json file to include the weights and preferred providers you'd like to use. 
        1. Details here: [rating-config.json.md](rating-config.json.md)

## Build the proxy-router Docker Image:
1. Build the Docker image using the provided `docker_build.sh` script
    * `./docker_build.sh --build` 
    * This script pulls the current Git tag and version, builds the docker image with all the port and defaults defined in the docker-compose.yml file
    
## Running the proxy-router Docker Container:
1. Run the Docker container using the provided `docker_build.sh` script
    * `./docker_build.sh --run`
    * This script runs the docker container with the port and volume mappings defined in the docker-compose.yml file
    * The proxy-router will start and begin listening for blockchain events

## Validating Steps:
1. Once the proxy-router is running, you can navigate to the Swagger API Interface (http://localhost:8082/swagger/index.html as example) to validate that the proxy-router is running and listening for blockchain events
1. Once validated, you can move on and create your provider, model and bid on the blockchain [03-provider-offer.md](03-provider-offer.md)

## NOTES:
* We have created docker-compose.yml, Dockerfile and docker_build.sh scripts to ease configuration and deployment of the proxy-router in a containerized environment
* Use these files as guides for applying to your system needs and configurations, private key, eth_node, ports, endpoints, volumes, .env, models-config.json and rating-config.json will need to be adjusted to your specific needs
* In most cases, the default .env file will work for the proxy-router...In some cases you will want to modify the .env file with advanced capability (log entries, private keys, private endpoints, etc)
* Please see the following for more information on these key config files: 
    * [proxy-router.all.env](proxy-router.all.env)
    * [models-config.json.md](models-config.json.md)
    * [rating-config.json.md](rating-config.json.md)