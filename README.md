# ACE* - Analytics of Covid Exposure Networks

![Alt text](ace-net.PNG?raw=true "ACE*")

## Quick links

- [About](#about)
- [Objectives](#objectives)
- [Background](#background)
- [ACE* Framework](#ace_framework)
- [Getting started](#getting_started)
- [Topics for future work](#future_work)
- [Interesting references](#references)


## About <a name="about"/>

The project explores possible applications of functional encryption (FE) on analyzing contact log data collected by [Corona-Warn-App](https://github.com/corona-warn-app) (CWA) on mobile devices, based on [Exposure Notification protocol](https://blog.google/documents/Exposure_Notification_-Cryptography_Specification_v1.2.1.pdf). 

Exposure Notification protocol collects contact log data through registration of Bluetooth low energy transmissions in a highly privacy preserving manner, on user's mobile devices, and prevents reconstruction of contact/exposure network, including the information about the contacts and their (possible) contagion. 

The goal of [ACE* project](https://github.com/cryptohackathon/ace-net) is to demonstrate proof-of-concept modifications of the current version of the CWA apps's protocol, where mobile devices could transmit certain functionally encrypted data derived from its contact tracing logs to a central analytics server, which would enable calculation of contact network parameters related to social distancing, while still fully preserving privacy.

In this way, users and health authorities could benefit from insights from the field about social distancing and aggregated exposure risk.


## Objectives <a name="objectives"/>

The objective of this project is to get insights about the actual social distancing during Covid-19 epidemia, directly from the participants.

Our plan is to leverage the data of actual exposure networks, gathered by multiple Corona-Warn-App (CWA) apps, and their analysis, using [decentralized multi- client functional encryption for inner product ](https://eprint.iacr.org/2017/989.pdf) FE scheme, while at all steps respecting data privacy and safety, by processing a minimum of required personal data that is handled with maximum protection.

We have implemented a proof-of-concept (PoC) implementation of the core infrastructure for trusted submitting, processing, and visualization of Covid-19 exposure networks, together with working simulation. 

Next step should be to bring these new functionalities to the CWA app.


## Background <a name="background"/>

[Corona-Warn-App (CWA)](https://www.coronawarn.app) is an open source project (mobile app + servers) that helps in tracing infection chains of SARS-CoV-2. It is uses a decentralized approach, with focus on data privacy and safety, (see e.g. privacy-preserving contact tracing specifications by [Apple](https://covid19.apple.com/contacttracing) & [Google](https://www.google.com/covid19/exposurenotifications)), and notifies users if they have been exposed. CWA is specifically designed to [ensure](https://blog.google/documents/Exposure_Notification_-Cryptography_Specification_v1.2.1.pdf) for each step that the app processes a minimum of required personal data that is handled with maximum protection - for the following 2 objectives:

1. Assess personal risk of infection, where the app
    - automatically collects nearby identifiers,
    - distributes list of keys of SARS-CoV-2 confirmed users,
    - checks for exposure to SARS-CoV-2 confirmed users,

2. Learn COVID-19 test results faster, where the app
    - enables communication (retrieve, inform) of user’s test result, after explicit consent from the user.   

In the context of exposure network analytics, where privacy of a user and his/her sensitive/health data (risk profile, encounters, data related to social distancing, etc.) is of utmost importance, decentralized functional encryption schemes can be used to securely calculate operations on the top of the encrypted data about the structure of exposure networks, comming from multiple users. For example, [GoFE](https://github.com/fentec-project/gofe) and [CiFEr](https://github.com/fentec-project/CiFEr) APIs provide cryptographic schemes for decentralization and functionalities up to quadratic functions. 

The communication between users's device (equipped with CWA app) and the server (CWA Infrastructure) is automated, mainly using pull mechanism, initiated from the user's device. The only exception is when users submits the positive test result for all user's encounters to be notified about possible infection. Mobile devices communicate only with servers and randomly submit requests, which are ignored by the servers, to prevent the attack from the malicious third party monitoring the traffic looking for communication patterns 

Hence, calculations and possible exposure checking always happen on the device; e.g. the device once per day pulls the data of all diagnosis keys (a.k.a temporary exposure keys of contagious people) and compares them with the encounter keys stored on the device from the given time interval. At these times, some interesting properties of the local encounter network are calculated and made available to the user, see e.g. [risk score calculation](https://www.r-bloggers.com/2020/09/risk-scoring-in-digital-contact-tracing-apps).

In the current version of the protocol, the data about the network and all encounters are stored in a secure storage on the device, see e.g. [here](https://covid19.apple.com/contacttracing). Due to privacy preservation, nobody (including the app) has access to the network, including the local encounters. Hence, (without the protocol upgrade) we cannot reconstructure the parts of the network and can only use the data available on the application level, specifically:

- The number of days since the most recent exposure.
- For up to last 14 days:
    - The number of keys that matched for an exposure detection.
    - The highest, full-range risk score of all user’s exposures.
    - The highest risk score of all exposure incidents.
    - The sum of the full-range risk scores for all exposures for the user.

Actual information from the field, e.g. about social distancing (days since exposure), yesterday’s risk encounters (nr. of keys matched), and yesterday’s risk received (max, sum), can help in better management of social distancing and the health crysis. Health authorities now have additional tool to better specify the configuration settings that are used in risk score calculations and notification thresholds, enabling them to control the virus, by using less coercive measures. Additionally, by lowering risky encounters (e.g. through efficient policies), the spread of the disease can be lowered, too.

Each CWA app can share (after explicit confirmation by the user) data about user's risk exposures or data related to social distancing to analytics server, that calculates totals from multiple participants and prepares the visualisations with meaningful insights about exposure networks. 
For "days since exposure" we create a vector of lenght 14 (one bucket for each day), whereas for "the number of keys that matched for an exposure detection" we create a vector of length seven (with buckets  0/1-2/3-5/6-10/11-20/21-40/41+, each representing an interval for the number of exposures matched). These data can be processed for a region, when region label is provided by the user or for the whole country. We use functional encryption (i.e. decentralized multi-client functional encryption scheme for inner product)  to encrypt and process the data from multiple CWA apps on the analytics server, see bellow.



## ACE* Framework <a name="ace_framework"/>

In order to gather encrypted data and extract some useful metrics from it, we have developed a framework with the decentralized multi-client functional encryption scheme for inner product ([DMCFE](https://eprint.iacr.org/2017/989.pdf)) in its core. The framework serves to produce data represented in histograms. A histogram is obtained from the data submitted by a group of clients (CWA users) who join to a pool characterized by time and location.

![Alt text](ace-net-scheme.png?raw=true "ACE* Framework")

### Process
* *Phase 1: Registration*. Clients join a pool. On enrollment they are assigned a sequence number. Based on the sequence number a client generates a key and submits it to the pool.

* *Phase 2: Key sharing*. All client keys are gathered in a pool and communicated to the clients together with histogram details.

* *Phase 3: Data encryption and key derivation*. Using the agreed multiple key, each client encripts its data (usually a vector of zeros with a single one corresponding to the bin of the histogram) and derives a key share (usually a key share for a vector of ones). The client sends this data back to pool.

* *Phase 4: Data decryption*. All data is gathered in the pool. It consist of a ciphertext collected from clients and key shares than can be used to apply an inner product operation to the ciphertexts. This in turn allows a histogram extraction in the analytics server.

### Components

* *Backend app*. A node.js application for pool management. It takes care of pool initialization and exposes an API to clients. It stores the data received from clients and eventually communicates with analytics server to extract histograms.

* *Client CLI in the client mode*. A command line interface written in GO is a user app mockup that communicates with the backend app. It is based on the [GoFE](https://github.com/fentec-project/gofe) library and uses DMCFE scheme.

* *Client CLI in the analytics mode*. A command line interface written in GO is an analytics server mockup that is called by the backend app to decrypt user data. It is based on the [GoFE](https://github.com/fentec-project/gofe) library.

* *Frontend app*. A node.js application build with Angular CLI that interactively logs the data gathering process.


## Getting started <a name="getting_started"/>

This is a quick intro to run demo. Find more detailed instructions in the project subdirectories.

### [api-back](/api-back)
Required: [node.js](https://nodejs.org/)

Run `npm install` and `npm run dev` to start the server on port 9500. Access Swagger with [API documentation](https://cryptohackathon.github.io/ace-net/) on `localhost:9500/api-doc/`.

### [ace-net-fe](ace-net-fe)
Required: [node.js](https://nodejs.org/), [Angular CLI](https://github.com/angular/angular-cli)

Run `npm install` and `ng serve` to start the server on port 4200.

### [test-clients](/test-clients)
Required: [go](https://golang.org/), [GoFE](https://github.com/fentec-project/gofe)

Use `./run-clients.sh`to start client simulator and `./run-analytics-server` to start mockup analytics server.


## Topics for future work <a name="future_work"/>

- Bring the results of this project to Corona-Warn-App.
- “Those who count the votes decide everything” - Voting systems (giving N votes between M >= N options).
- The decentralized nature of the framework fits perfectly to blockchain technologies. The pool management implemented in our backend system (api-back) could be transferred to e.g. [Ethereum Virtual Machine](https://ethereum.org/) ecosystem based on smart contracts.
- Randomising communication patterns, by involving a part (half?) of the requests to be ignored or by submitting data using secret sharing schemes (e.g. partial data, multiple times).
- GoFE problem: "panic: runtime error: invalid memory address or nil pointer dereference
 [signal SIGSEGV: segmentation violation code=0x1 addr=0x10 pc=0x67689b]".



## Interesting references <a name="references"/>

FE

- [Decentralized Multi-Client Functional Encryption for Inner Product](https://eprint.iacr.org/2017/989.pdf)
- Subway heatmap [paper](https://fentec.eu/content/privacy-enhanced-machine-learning-functional-encryption), [github repository](https://github.com/fentec-project/FE-anonymous-heatmap)

COVID

- [COVID-19 Tracker](https://covid-19.sledilnik.org)
- [Understanding metropolitan patterns of daily encounters](https://www.researchgate.net/publication/235009037_Understanding_metropolitan_patterns_of_daily_encounters)
- [Slovenian Temporary Exposure Keys](https://github.com/sledilnik/cwa-scrape)

Corona-Warn-APP

- Corona-Warn-App, [project](https://www.coronawarn.app), [github repository](https://github.com/corona-warn-app/cwa-documentation)
- [CWA Solution Architecture](https://github.com/corona-warn-app/cwa-documentation/blob/master/solution_architecture.md)
- Exposure notification framework for contact tracing, [Apple](https://covid19.apple.com/contacttracing)
[Google](https://www.google.com/covid19/exposurenotifications)
- [Exposure Notification cryptography specification](https://blog.google/documents/Exposure_Notification_-Cryptography_Specification_v1.2.1.pdf)
- [CWA Risk score calculation](https://www.r-bloggers.com/2020/09/risk-scoring-in-digital-contact-tracing-apps)

Contact tracing

- [Kukkala VB, Saini JS, Iyengar SRS, Privacy preserving network analysis of distributed social networks](https://eprint.iacr.org/2016/427.pdf)
- [Kukkala VB, Iyengar SRS, Computing Betweenness Centrality: An Efficient Privacy-Preserving Approach](https://link.springer.com/chapter/10.1007/978-3-030-00434-7_2)
- [Kukkala VB, Iyengar SRS, Identifying Influential Spreaders in a Social Network (While Preserving Privacy)](https://content.sciendo.com/downloadpdf/journals/popets/2020/2/article-p537.pdf)

ACE* project github repository

- [Github reposotory](https://github.com/cryptohackathon/ace-net)
- [API documenttion](https://cryptohackathon.github.io/ace-net)
