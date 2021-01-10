# ACE* - Analytics of Covid Exposure Networks

We will explore possible applications of functional encryption (FE) on analyzing contact log data collected by [Corona-Warn-App](https://github.com/corona-warn-app) on mobile devices, based on [Exposure Notification protocol](https://blog.google/documents/Exposure_Notification_-Cryptography_Specification_v1.2.1.pdf). 

Exposure Notification protocol collects contact log data through registration of Bluetooth low energy transmissions in a highly privacy preserving manner, on user's devices, that prevents reconstruction of contact/exposure network, including the information about the user and his/her (possible) contagion. 

Our goal is to propose possible modifications of the current version of the protocol where mobile devices could transmit certain functionally encrypted data from contact tracing logs to a central server (or between each other) which would enable calculation of certain contact network parameters while still fully preserving privacy. 

In this way, users and health authorities could benefit from information on a level of social distancing without invading their citizen's privacy.

## More details

[GoFE](https://github.com/fentec-project/gofe) and [CiFEr](https://github.com/fentec-project/CiFEr) APIs provide cryptographic schemes for decentralization and functionalities up to quadratic functions. 

In the context of exposure network analytics, where privacy of a user and his/her sensitive/health data is of utmost importance, decentralized functional encryption schemes can be used to securely calculate operations on the top of the encrypted data about the structure of exposure networks, comming from multiple users.

The communication between users's device (Corona Warn App) and the server (CWA Infrastructure) is using pull mechanism, initiated from the user's device. Hence calculations currently always happen on the (mobile) device; e.g. the device once per day pulls the data of all diagnosis keys (a.k.a temporary exposure keys of contagious people) and compares them with the encounter keys stored on the device from the given time interval. At these times, the structure or interesting properties of the local encounter network could be calculated and encrypted.

Randomly, to prevent the attack from the malicious third party monitoring the traffic looking for communication patterns, mobile device communicates with servers and submits requests that will be ignored by the servers. These requests could be used to submit encrypted data about the structure or properties of the device's encounter network.


## Interesting topics

- Uploading temporary exposure keys (TEKs) after a user submits a positive test (the user is still considered contagious for few days) through exposure network
- Calculation: size of my exposure network, nr. of exposures, max cluster, max clique size


## Selected challenge = Apply

Demonstrate a viable application of FE, for example, in the use case of access control or privacy-preserving AI, see [challenges](https://cryptohackathon.eu/#challenges) and [tutorial](https://us02web.zoom.us/rec/share/PeSRUAaUYbDBiN6AaQeszotTeALfuDyMyZxX5TbnfQxaGUGl4H_DGOCMeUmEToMD.0x18XGtN1pHHKGb1) (password: W$&k7Gwg).

Show the practicability and versatility of FE and reduce the entry barrier for researchers, engineers, product owners, architects to deploy FE in their solutions, ranging from prototypes to enterprise-ready products.

- Technical Execution (40 %)
    - How has the team effectively utilized the FE technologies?
    - How easy is the application to use?
    - How advanced is the prototype presented?
    - Is there a working demo?
    - Is the project well documented (readme, wiki), such that follow-up work is possible.
- Challenge Fit (20 %)
    - How relevant is the project to the stated challenge?
- Impact (40 %)
    - Will this solution have a far reach and market potential?



## Interesting references

- [Exposure Notification cryptography specification](https://blog.google/documents/Exposure_Notification_-Cryptography_Specification_v1.2.1.pdf)
- [Corona-Warn-App documentation](https://github.com/corona-warn-app/cwa-documentation)
- [Corona-Warn-App Solution Architecture](https://github.com/corona-warn-app/cwa-documentation/blob/master/solution_architecture.md)
- Kukkala VB, Saini JS, Iyengar SRS, [Privacy preserving network analysis of distributed social networks](https://eprint.iacr.org/2016/427.pdf), International Conference on Information Systems Security, 336-355, 2016
- VB Kukkala, SRS Iyengar, [Computing Betweenness Centrality: An Efficient Privacy-Preserving Approach](https://link.springer.com/chapter/10.1007/978-3-030-00434-7_2), International Conference on Cryptology and Network Security, 23-42
- VB Kukkala, SRS Iyengar, [Identifying Influential Spreaders in a Social Network (While Preserving Privacy)](https://content.sciendo.com/downloadpdf/journals/popets/2020/2/article-p537.pdf), Proceedings on Privacy Enhancing Technologies 2020 (2), 537-557


