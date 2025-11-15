package main

import (
	"context"
	"log"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		r, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			return nil, err
		}
		return []*Account{{
			ID:   r.ID,
			Name: r.Name,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}
	as, err := r.server.accountClient.ListAccounts(ctx, skip, take)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var as_ []*Account
	for _, v := range as {
		as_ = append(as_, &Account{
			ID:   v.ID,
			Name: v.Name,
		})
	}
	return as_, nil
}

func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		r, err := r.server.catalogClient.GetProduct(ctx, *id)
		if err != nil {
			return nil, err
		}
		return []*Product{{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Price:       r.Price,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}
	q := ""
	if query != nil {
		q = *query
	}
	ps, err := r.server.catalogClient.GetProducts(ctx, skip, take, nil, q)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var ps_ []*Product
	for _, v := range ps {
		ps_ = append(ps_, &Product{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
		})
	}
	return ps_, nil
}

func (p *PaginationInput) bounds() (uint64, uint64) {
	s_ := uint64(0)
	t_ := uint64(10)
	if p.Skip != nil {
		s_ = uint64(*p.Skip)
	}
	if p.Take != nil {
		t_ = uint64(*p.Take)
	}

	return s_, t_
}
