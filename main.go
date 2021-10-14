package main

import (
	"context"
	"log"
	"math/big"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type airdrop struct {
	Address  string `json:"address"`
	Quantity int    `json:"quantity"`
}

var client *rpc.Client

func createAirdrop(c echo.Context) error {

	addr := &airdrop{}

	if err := c.Bind(addr); err != nil {
		return err
	}

	pubKey := solana.MustPublicKeyFromBase58(addr.Address)

	// Airdrop 100 SOL to the new account:
	out, err := client.RequestAirdrop(
		context.TODO(),
		pubKey,
		solana.LAMPORTS_PER_SOL*uint64(addr.Quantity),
		rpc.CommitmentFinalized,
	)

	if err != nil {
		log.Println(err)
	}

	log.Println("airdrop transaction signature:", out)

	return c.JSON(http.StatusCreated, out)
}

func getBalance(c echo.Context) error {

	addr := c.Param("address")

	pubKey := solana.MustPublicKeyFromBase58(addr)

	out, err := client.GetBalance(
		context.TODO(),
		pubKey,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		log.Println(err)
	}

	spew.Dump(out)
	spew.Dump(out.Value) // total lamports on the account; 1 sol = 1000000000 lamports

	var lamportsOnAccount = new(big.Float).SetUint64(uint64(out.Value))
	// Convert lamports to sol:
	var solBalance = new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))

	// WARNING: this is not a precise conversion.
	return c.JSON(http.StatusOK, solBalance.Text('f', 10))
}

func main() {

	// Create a new RPC client:
	client = rpc.New(rpc.TestNet_RPC)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/airdrop", createAirdrop)
	e.GET("/address/:address", getBalance)

	// Start server
	e.Logger.Fatal(e.Start(":1337"))
}
