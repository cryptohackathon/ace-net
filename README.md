# ACE* - Analytics of Covid Exposure Networks

We will explore possible applications of functional encryption (FE) on analyzing contact log data collected by [Corona-Warn-App](https://github.com/corona-warn-app) on mobile devices, based on [Exposure Notification protocol](https://blog.google/documents/Exposure_Notification_-Cryptography_Specification_v1.2.1.pdf). Exposure Notification protocol collects contact log data through registration of Bluetooth low energy transmissions in a highly privacy preserving manner that prevents reconstruction of contact/exposure network. Our goal is to propose possible modifications of the protocol where mobile devices could transmit certain functionally encrypted data from contact tracing logs to a central server (or between each other) which would enable calculation of certain contact network parameters while fully preserving privacy. In this way, authorities could benefit from information on level of social distancing without invading their citizen's privacy.

## Interesting references

- [Exposure Notification cryptography specification](https://blog.google/documents/Exposure_Notification_-Cryptography_Specification_v1.2.1.pdf)
- [Corona-Warn-App documentation](https://github.com/corona-warn-app/cwa-documentation)
- Kukkala VB, Saini JS, Iyengar SRS, [Privacy preserving network analysis of distributed social networks](https://eprint.iacr.org/2016/427.pdf), International Conference on Information Systems Security, 336-355, 2016
- VB Kukkala, SRS Iyengar, [Computing Betweenness Centrality: An Efficient Privacy-Preserving Approach](https://link.springer.com/chapter/10.1007/978-3-030-00434-7_2), International Conference on Cryptology and Network Security, 23-42
- VB Kukkala, SRS Iyengar, [Identifying Influential Spreaders in a Social Network (While Preserving Privacy)](https://content.sciendo.com/downloadpdf/journals/popets/2020/2/article-p537.pdf), Proceedings on Privacy Enhancing Technologies 2020 (2), 537-557


