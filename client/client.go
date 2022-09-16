package client

import (
	"context"
	"log"

	pb "github.com/DonggyuLim/grc20/protos/RPC"
	u "github.com/DonggyuLim/uniswap/utils"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
)

type Client struct {
	client pb.RPCClient
}

var (
	cli *Client
	ctx context.Context
)

const serviceHost = "localhost:9000"

func GetRPCClient() *Client {

	if cli == nil {
		conn, err := grpc.Dial(serviceHost, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		cli = &Client{
			client: pb.NewRPCClient(conn),
		}
	}
	return cli
}

func (c *Client) GetBalance(tokenName, account string) (decimal.Decimal, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := c.client.GetBalance(ctx, &pb.GetBalanceRequest{
		TokenName: tokenName,
		Account:   account,
	})
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	return u.NewDecimalFromUint(res.GetBalance()), nil
}

func (c *Client) TokenInfo(tokenName string) (*pb.TokenInfoResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := c.client.GetTokenInfo(ctx, &pb.TokenInfoRequest{
		TokenName: tokenName,
	})
	if err != nil {
		return &pb.TokenInfoResponse{}, err
	}
	return res, nil
}

func (c *Client) Allowance(tokenName, owner, spender string) (decimal.Decimal, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := c.client.GetAllowance(ctx, &pb.AllowanceRequest{
		TokenName: tokenName,
		Owner:     owner,
		Spender:   spender,
	})
	if err != nil {
		return u.NewDecimalFromUint(0), err
	}
	return u.NewDecimalFromUint(res.GetAllowance()), nil
}

func (c *Client) Approve(tokenName, owner, spender string, amount uint64) (*pb.ApproveResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := c.client.Approve(ctx, &pb.ApproveRequest{
		TokenName: tokenName,
		Owner:     owner,
		Spender:   spender,
		Amount:    amount,
	})
	if err != nil {
		return &pb.ApproveResponse{}, err
	}
	return res, nil
}

func (c *Client) TransferFrom(tokenName, from, to, spender string, amount uint64) (*pb.TransferFromResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := c.client.TransferFrom(ctx, &pb.TransferFromRequest{
		TokenName: tokenName,
		From:      from,
		To:        to,
		Spender:   spender,
		Amount:    amount,
	})
	if err != nil {
		return &pb.TransferFromResponse{}, err
	}
	return res, nil
}

func (c *Client) Transfer(tokenName, to, from string, amount uint64) (*pb.TransferResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := c.client.Transfer(ctx, &pb.TransferRequest{
		TokenName: tokenName,
		To:        to,
		From:      from,
		Amount:    amount,
	})
	if err != nil {
		return &pb.TransferResponse{}, err
	}
	return res, nil
}
