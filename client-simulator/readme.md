# Client Simulator

Simulator is based on the Heatmap example ([fentec-project/FE-anonymous-heatmap](https://github.com/fentec-project/FE-anonymous-heatmap)) and uses the decentralized multi-client functional encryption scheme for inner product ([paper](https://eprint.iacr.org/2017/989.pdf)) implementation in [fentec-project/gofe](https://github.com/fentec-project/gofe). It simulates a participation of clients in pools that allows an anonymous extraction of data histograms. Clients in a pool reach a consensus on secret keys and then share encrypted data. This eventually allows the analytics engine to produce histograms of client submitted data.

To run the simulator, download [fentec-project/gofe](https://github.com/fentec-project/gofe):
```
go get -u -t github.com/fentec-project/gofe/...
```
and execute
```
go run .
```
The histograms resulting from the simulation are stored in text files in the subfolder `/data`.