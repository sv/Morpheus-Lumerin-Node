# Morpheus Lumerin Node
The purpose of this software is to enable interaction with distributed, decentralized LLMs on the Morpheus network through a desktop chat experience.

The steps listed below are for both the Consumer and Provider to get started with the Morpheus Lumerin Node. As the software is developed, both onboarding & configuration of the provider and consumer roles will be simplified, automated and more transparent to the end user.

# **NOTE: ARBITRUM SEPOLIA TESTNET ONLY at this time - DEVELOPER PREVIEW ONLY**

**Components that are included in this repository are:**
* Lumerin Proxy-Router is a background process that monitors sepcific blockchain contract events, manages secure sessions between consumers and providers and routes prompts and responses between them
* Lumerin UI-Desktop is the front end UI to interact with LLMs and the Morpheus network via the proxy-router
* Local Llama.cpp and tinyllama model to run locally 

## Tokens and Contract Information 
* Morpheus saMOR Token: `0xc1664f994fd3991f98ae944bc16b9aed673ef5fd` 
* Lumerin Morpheus Smart Contract : `0x8e19288d908b2d9F8D7C539c74C899808AC3dE45`
    * Interact with the Morpheus Contract: https://louper.dev/diamond/0x8e19288d908b2d9F8D7C539c74C899808AC3dE45?network=arbitrumSepolia#write
* Blockchain Explorer: `https://sepolia.arbiscan.io/`
* Swagger API: `http://localhost:8082/swagger/index.html`

## Prerequisites
* **WALLET:** For testing, you will need both `saMOR` and `saETH` tokens in your wallet. You should be able to get either of these from the usual Sepolia Arbitrum testnet faucets.
    * `saMOR` is the token used to pay for the model provider staking and consumer usage
    * `saETH` is the token used to pay for the gas fees on the network  
    * At this time, we recommend setting up a new wallet in either MetaMask or other existing crypto wallet 
        * You will need both the wallet's `private key` and `mnemonic` from your wallet to startup and interact with both the proxy-router and the UI-Desktop 
        * AS ALWAYS:  **NEVER** share either of these two items with anyone else, as they are the keys to your wallet and can be used to access your funds
        * In the future, as the UI-Desktop functionality is developed, all of this will be handled in the UI and you will not need to manually enter these items

* **ETH NODE:** You will need access to either a public or private Sepolia Arbitrum ETH node that the proxy-router and ui-desktop will use to monitor events on the blockchain
    * We have provided an existing public node example `https://arbitrum-sepolia.blockpi.network/v1/rpc/public` in the .env-example file for you to use, however, 
    * We stronly recommend either a wss: or https: private endpoint from a trusted service like Alchemy or Infura (both offer free accounts for testing/private usage)


## Getting Started (using the provided releases - run from source instructions coming soon)
### Common Steps for both Consumer and Provider: 
1. PreRequisites: 
    1. Have Mnemonic and Private key for a Consumer Wallet (Suggest setting one up in MetaMask so you can pull the mnemonic and private key) 
    2. Obtain or Transfer saMOR and saETH to this wallet for testing (suggest 1000 saMOR and at least 0.1 saETH) - Sepolia Arbitrum Chain 
2. Download latest release for your operating system: https://github.com/Lumerin-protocol/Morpheus-Lumerin-Node/releases
3. Extract the zip to a local folder (examples)
    * Windows: `(%USERPROFILE%)/Downloads/morpheus)` 
    * Linux & MacOS: `~/Downloads/morpheus`
    * On MacOS you may need to execute `xattr -c mor-launch proxy-router ui-desktop.app llama-server` in a command window to remove the quarantine flag on MacOS
4. Edit the .env file (this is a hidden file, please use your OS specific method to show hidden files) 
    * Change `ETH_NODE_ADDRESS=` (you can setup a free one in Alchemy or Infura) 
    * Change `WALLET_PRIVATE_KEY=` This will be the private key of the Wallet you setup previously

### Consumer from Package (Local LLM, Proxy-Router & UI-Desktop): 
Please see /consumer.md for building from git source 
1. Assuming that you have already setup the prerequisites and downloaded the latest release, you can follow the steps below to get started
2. Launch the node - this should open a command window to see local LLM model server and proxy-router start and then should launch the user interface  
    * Windows: Double click the `mor-launch.exe` (You will need to tell Windows Defender this is ok to run) 
    * Linux & MacOS: Open a terminal and navigate to the folder and run `./mor-launch`
3. Startup User Interface: 
    1. Read & accept terms & Conditions 
    2. Set a strong password (this is for the UI-Desktop only)
    3. When prompted to create the wallet from a new Mnemonic, select **Recover your wallet** from `Saved Mnemonic` instead.  
        * This is important so that the private key for the proxy-router (in the .env file) and the wallet running in the UI are the same 
        * If you create a new wallet, the proxy-router will be listening to a different wallet address than the UI is managing and this will not work.
4. Local Test: Once the UI is up and running,
    1. You should see tokens for saETH and saMOR that you sent to this wallet earlier. If not, either the ETH_NODE_ADDRESS is incorrect or the wallet address is not the same as the one you sent the tokens to
    2. Click on the `Chat` icon on the left side of the screen
    3. Make sure the `Local Model` is selected
    4. Begin your conversation with the model by typing in the chat window and pressing `Enter`
        * You should see the model respond with the appropriate response to the prompt you entered, if not, then likely the configuration in the .env file is incorrect or the model is not running correctly
5. Remote Test: Once you've verified that your wallet can access the blockchain and you can see the local model working, you can switch to a remote model and test that as well
    1. In the `Chat` window, select `Change Model `
        1. Select a different model from remote providers
        2. DropDown and select the contract address of the model you want to use 
        3. Click Change 
        4. Click Open Session 
        5. MUST Enter at least **5** MOR to open session 
    3. You can now chat with the remote model and see the responses in the chat window 
6. Cleanup/Closeout 
    * Manually End all Remote Sessions: 
        * In the Chat Window, click on the Time icon to the right of the Model line - this will expand and show current sessions, click the "X" next to each one to make sure it's closed 
    * Closing the UI-Desktop window should leave the CMD window open
        * You’ll have to ctrl-c in the window to kill the local model and proxy-router
    * To FULLY delete and force a clean startup of the UI (including forcing new password and mnemonic recovery), delete the ui-desktop folder and start the UI again
        * Windows:  `%USERPROFILE%\AppData\Roaming\ui-desktop`
        * Linux: `~/.config/ui-desktop`
        * MacOS: `~/Library/Application Support/ui-desktop`

### Provider (Local LLM to offer, Proxy-Router running as background/service): 
This section is used for offering your hosted LLM model to the network for others to use.
1. LAUNCH LOCAL MODEL:
    * On your server launch the local model server
    * You can use the provided `llama.cpp` and `tinyllama` model to test locally
    * If your local model is listening on a different port locally, you will need to modify the `OPENAI_BASE_URL` in the .env file to match the correct port
1. LAUNCH PROXY-ROUTER: 
    * **Your proxy-router must have a publicly accessible endpoint for the provider (ip:port or fqdn:port no protocol) eg: `mycoolmornode.domain.com:3989` - this will be used when creating the provider on the blockchain** 
    * On your server, launch the proxy-router with the modified .env file shown in the common pre-requisites section
    * Windows: Double click the `proxy-router.exe` (You will need to tell Windows Defender this is ok to run)  
    * Linux & MacOS: Open a terminal and navigate to the folder and run `./proxy-router`from the morpheus/proxy-router folder
    * This will start the proxy-router in the background and begin monitoring the blockchain for events

1. **(OPTIONAL)** - External Provider or Pass through 
    * In some cases you will want to leverage external or existing AI Providers in the network via their own, private API
    * Dependencies: 
        * `model-config.json` file in the proxy-router directory
        * proxy-router .env file for proxy-router must also be updated to include `MODELS_CONFIG_PATH=<path_to_proxy-router>/models-config.json`
    * Once your provider is up and running, deploy a new model and model bid via the diamond contract (you will need the `model_ID` for the configuration)
    * Edit the model-config.json to the following json format ... with 
    * The JSON ID will be the ModelID that you created above, modelName, apiTYpe, apiURL and apiKey are from the external provider and specific to their offered models 
    * Once the model-config.json file is updated, the morpheus node will need to be restarted to pick up the new configuration (not all models (eg: image generation can be utilized via the UI-Desktop, but API integration is possible)
1. SETUP PROVIDER / MODEL / BID: 
    1. WEB3/Arbiscan/Metamask: Authorize Diamond Contract to spend on the Provider's behalf 
       1. https://sepolia.arbiscan.io/address/0xc1664f994fd3991f98ae944bc16b9aed673ef5fd#writeContract 
       1. Connect to Web3 (connect Provider wallet) 
       1. Click Approve 
       1. Spender Address = Diamond Contract 
       1. Authorized Amount = remember that this is in the form 1*10^18 so make sure there's enough MOR 1ranted to cover the contract fees 
       1. The Diamond Contract is now authorized to spend MOR on provider's behalf 
    1. Create Provider in the Diamond contract via swagger api:
        1. Start proxy-router 
        1. http://localhost:8082/swagger/index.html#/providers/post_blockchain_providers
        1. Enter required fields:   
            1. addStake = Amount of stake for provider to risk - Stake can be 0 now 
            1. Endpoint = Your publicly accessible endpoint for the proxy-router provider (ip:port or fqdn:port no protocol) eg: `mycoolmornode.domain.com:3989`
    1. Create Model in the contract:
        1. Go to http://localhost:8082/swagger/index.html#/models/post_blockchain_models and enter
            1. modelId: random 32byte/hex that will uniquely identify model (uuid)
            1. ipfsCID: another random32byte/hex for future use (model library)
            1. Fee: fee for the model usage - 0 for now
            1. addStake: stake for model usage - 0 for now 
            1. Owner: Provider Wallet Address 
            1. name: Human Readable model like "Llama 2.0" or "Mistral 2.5" or "Collective Cognition 1.1" 
            1. tags: array of tag strings for the model 
            1. Capture the `modelID` from the JSON response
    4. Offer Model Bid in the contract: 
        1. Navigate to http://localhost:8082/swagger/index.html#/bids/post_blockchain_bids and enter
            1. modelID: Model ID Created in last step: 
            1. pricePerSecond: this is in 1*10^18 format so 100000000000 should make 5 minutes for the session 1round 37.5 saMOR 
            1. Click Execute 


```
#model-config.json 
{
    "0x4b5d6c2d3e4f5a6b7c8de7f89a0b19e07f4a6e1f2c3a3c28d9d5e6": {
        "modelName": "v1-5-specialmodel.modelversion [externalmodel]",
        "apiType": "provider_api_type",
        "apiUrl": "https://api.externalmodel.com/v1/xyz/generate",
        "apiKey": "api-key-from-external-provider"
    },
    "0xb2c8a6b2c1d9ed7f0e9a3b4c2d6e5f14f9b8c3a7e5d6a1a0b9c7d8e4f30f4a7b": {
        "modelName": "v1-7-specialmodel2.modelversion [externalmodel]",
        "apiType": "provider_api_type",
        "apiUrl": "https://api.externalmodel.com/v1/abc/generate",
        "apiKey": "api-key-from-external-provider"
    }
}
```

5. VERIFY PROVIDER SETUP 
    * On a separate machine and with a separate wallet, you can follow the consumer steps above to verify that your model is available and working correctly