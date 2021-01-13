/*
Hacked as part of crypto hackhaton, based on a sample FE-anonymous-heatmap (see https://github.com/fentec-project/FE-anonymous-heatmap)
*/

package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/fentec-project/bn256"
	"github.com/fentec-project/gofe/data"
	"github.com/fentec-project/gofe/innerprod/fullysec"
	"github.com/pkg/errors"
)

func main() {
	simulation()
}

func simulation() {
	numOfPools := 5
	numOfClientsPerPool := 10

	labels := []string{"0", "5", "10", "15", "20"}

	oneVec := data.NewConstantVector(numOfClientsPerPool, big.NewInt(1))
	vectors := make([]data.Vector, 1)
	vectors[0] = oneVec

	path := "data"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	pools := make([]*Pool, numOfPools)
	for i := 0; i < numOfPools; i++ {
		pools[i] = new(Pool)
		pools[i].init(fmt.Sprintf("histogram_%d", i), labels, vectors, numOfClientsPerPool)
		fmt.Printf("Pool %s: opened, awaiting %d clients\n", pools[i].code, numOfClientsPerPool)
	}

	p := rand.Perm(numOfPools)
	for i := 0; i < numOfPools; i++ {
		for j := 0; j < numOfClientsPerPool; j++ {
			go simulateClient(pools[p[i]])
		}
	}
	for i := 0; i < numOfPools; i++ {
		go simulateAnalyticsEngine(pools[p[i]])
	}

	time.Sleep(1000 * time.Second)
}

func simulateClient(pool *Pool) {
	client := new(Client)

	name := fmt.Sprintf("%d", 1000+rand.Intn(1000))

	dataLength := 5
	data := data.NewConstantVector(dataLength, big.NewInt(0))
	data[rand.Intn(dataLength)] = big.NewInt(1)

	client.init(data)

	time.Sleep(time.Duration(rand.Intn(50)) * time.Second)
	client.connectToPool(pool)
	fmt.Printf("  Client %s joined pool %s\n", name, pool.code)

	if pool.publicKeysGathered() {
		fmt.Printf("Pool %s: clients agreed on secret keys\n", pool.code)
	}

	for !client.tryPostData() {
		time.Sleep(time.Duration(rand.Intn(20)) * time.Second)
	}
	fmt.Printf("  Client %s submitted encrypted data to pool %s\n", name, pool.code)
}

func simulateAnalyticsEngine(pool *Pool) []*big.Int {
	time.Sleep(time.Duration(10 * time.Second))

	for !pool.clientDataGathered() {
		time.Sleep(time.Duration(10 * time.Second))
	}
	fmt.Printf("Pool %v: client data gathered, ready for analysis \n", pool.code)

	histogram := pool.getValues(0)
	writeVecToFile(fmt.Sprintf("data/%s.txt", pool.code), pool.labels, histogram)
	fmt.Printf("Pool %s: histogram decrypted %v\n", pool.code, histogram)

	return histogram
}

// POOL: an entity representing a group of users that join to aggregate encrypted data

type Pool struct {
	code        string
	labels      []string
	noOfClients int
	clientIndex int
	publicKeys  []*bn256.G1
	ciphers     [][]*bn256.G1
	vectors     []data.Vector
	keyShares   [][]data.VectorG2
}

func (p *Pool) init(code string, labels []string, vectors []data.Vector, noOfClients int) {
	p.code = code
	p.labels = labels
	p.noOfClients = noOfClients
	p.clientIndex = 0
	p.publicKeys = make([]*bn256.G1, noOfClients)
	p.ciphers = make([][]*bn256.G1, len(labels))
	for i := 0; i < len(labels); i++ {
		p.ciphers[i] = make([]*bn256.G1, noOfClients)
	}
	p.vectors = vectors
	p.keyShares = make([][]data.VectorG2, len(vectors))
	for i := 0; i < len(vectors); i++ {
		p.keyShares[i] = make([]data.VectorG2, noOfClients)
	}
}

func (p *Pool) addClient() int {
	if p.clientIndex == p.noOfClients {
		return -1
	}
	index := p.clientIndex
	p.clientIndex++
	if p.clientIndex == p.noOfClients {
		// fmt.Printf("Pool %s: full, %d clients gathered\n", p.code, p.noOfClients)
	}
	return index
}

func (p *Pool) addClientPublicKey(index int, publicKey *bn256.G1) {
	p.publicKeys[index] = publicKey
	if p.publicKeysGathered() {
		// fmt.Printf("Pool %s: clients agreed on secret keys\n", p.code)
	}
}

func (p *Pool) getPublicKeys() []*bn256.G1 {
	if p.publicKeysGathered() {
		return p.publicKeys
	} else {
		return nil
	}
}

func (p *Pool) publicKeysGathered() bool {
	for i := 0; i < p.noOfClients; i++ {
		if p.publicKeys[i] == nil {
			return false
		}
	}
	return true
}

func (p *Pool) setClientData(index int, data []*bn256.G1, keyShares []data.VectorG2) {
	for i := 0; i < len(p.labels); i++ {
		p.ciphers[i][index] = data[i]
	}
	for i := 0; i < len(p.vectors); i++ {
		p.keyShares[i][index] = keyShares[i]
	}
	if p.clientDataGathered() {
		// fmt.Printf("Pool %v: client data gathared, ready for analysis \n", p.code)
	}
}

func (p *Pool) getValues(vectorIndex int) []*big.Int {
	var err error

	if !p.clientDataGathered() {
		return nil
	}

	y := p.vectors[vectorIndex]
	keyShares := p.keyShares[vectorIndex]

	values := make([]*big.Int, len(p.labels))
	for i := 0; i < len(p.labels); i++ {
		values[i], err = fullysec.DMCFEDecrypt(p.ciphers[i], keyShares, y, p.labels[i], big.NewInt(int64(p.noOfClients)))
		if err != nil {
			panic(errors.Wrap(err, "could not decrypt"))
		}
	}

	return values
}

func (p *Pool) clientDataGathered() bool {
	for i := 0; i < p.noOfClients; i++ {
		if p.ciphers[0][i] == nil {
			return false
		}
	}
	return true
}

// CLIENT: an entity representing a person that cooperates in a pool data collection

type Client struct {
	pool   *Pool
	index  int
	client *fullysec.DMCFEClient
	data   data.Vector
}

func (c *Client) init(data data.Vector) {
	c.data = data
}

func (c *Client) connectToPool(pool *Pool) {
	var err error
	c.pool = pool
	c.index = pool.addClient()
	c.client, err = fullysec.NewDMCFEClient(c.index)
	if err != nil {
		panic(errors.Wrap(err, "could not instantiate fullysec"))
	}
	pool.addClientPublicKey(c.index, c.client.ClientPubKey)
}

func (c *Client) tryPostData() bool {
	publicKeys := c.pool.getPublicKeys()
	if publicKeys == nil {
		return false
	}
	err := c.client.SetShare(publicKeys)
	if err != nil {
		panic(errors.Wrap(err, "could not create private values"))
	}
	c.pool.setClientData(c.index, c.getDataEncrypton(), c.getKeyShares())
	return true
}

func (c *Client) getDataEncrypton() []*bn256.G1 {
	ciphers := make([]*bn256.G1, len(c.pool.labels))
	for i := 0; i < len(c.pool.labels); i++ {
		cipher, err := c.client.Encrypt(c.data[i], c.pool.labels[i])
		if err != nil {
			panic(errors.Wrap(err, "could not encrypt"))
		}
		ciphers[i] = cipher
	}
	return ciphers
}

func (c *Client) getKeyShares() []data.VectorG2 {
	keyShares := make([]data.VectorG2, len(c.pool.vectors))
	for i := 0; i < len(c.pool.vectors); i++ {
		keyShare, err := c.client.DeriveKeyShare(c.pool.vectors[i])
		if err != nil {
			panic(errors.Wrap(err, "could not generate key shares"))
		}
		keyShares[i] = keyShare
	}
	return keyShares
}
