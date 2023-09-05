package loadtest

import (
	"encoding/csv"
	"fmt"
	"os"
)

// ReapappClientFactory creates load testing clients to interact with the
// Cosmos SDK built-in simapp ABCI application.
type ReapappClientFactory struct{}

// ReapappClient is responsible for generating transactions. Only one client
// will be created per connection to the remote Tendermint RPC endpoint, and
// each client will be responsible for maintaining its own state in a
// thread-safe manner.
type ReapappClient struct {
	txs []string
	cnt int
}

// ReapapClientFactory implements loadtest.ClientFactory
var (
	_ ClientFactory = (*ReapappClientFactory)(nil)
	_ Client        = (*ReapappClient)(nil)
)

func init() {
	if err := RegisterClientFactory("reapapp", NewReapappClientFactory()); err != nil {
		panic(err)
	}
}

func NewReapappClientFactory() *ReapappClientFactory {
	return &ReapappClientFactory{}
}

func (f *ReapappClientFactory) ValidateConfig(cfg Config) error {
	maxTxsPerEndpoint := cfg.MaxTxsPerEndpoint()
	if maxTxsPerEndpoint < 1 {
		return fmt.Errorf("cannot calculate an appropriate maximum number of transactions per endpoint (got %d)", maxTxsPerEndpoint)
	}
	return nil
}

func (f *ReapappClientFactory) NewClient(cfg Config) (Client, error) {
	// Open the CSV file
	file, err := os.Open(cfg.TxInputFile)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return &ReapappClient{}, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading the CSV file:", err)
		return &ReapappClient{}, err
	}

	// Process the records (in this example, we just print them)
	txs := []string{}
	for i := 0; i < len(records); i++ {
		fmt.Println(records[i][0])
		txs = append(txs, records[i][0])
	}

	return &ReapappClient{
		txs: txs,
		cnt: 0,
	}, nil
}

func (c *ReapappClient) GenerateTx() (string, error) {
	tx := c.txs[c.cnt]
	c.cnt += 1
	return tx, nil
}
