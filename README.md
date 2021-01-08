# ACE* - Analytics of Covid Exposure Networks

We will explore possible applications of functional encryption (FE) on analyzing contact log data collected by (https://github.com/corona-warn-app)[Corona-Warn-App] on mobile devices, based on (https://blog.google/documents/Exposure_Notification_-Cryptography_Specification_v1.2.1.pdf)[Exposure Notification protocol]. Exposure Notification protocol collects contact log data through registration of Bluetooth low energy transmissions in a highly privacy preserving manner that prevents reconstruction of contact/exposure network. Our goal is to propose possible modifications of the protocol where mobile devices could transmit certain functionally encrypted data from contact tracing logs to a central server (or between each other) which would enable calculation of certain contact network parameters while fully preserving privacy. In this way, authorities could benefit from information on level of social distancing without invading their citizen's privacy.

## Interesting references

- (https://blog.google/documents/Exposure_Notification_-Cryptography_Specification_v1.2.1.pdf)[Exposure Notification cryptography specification]

