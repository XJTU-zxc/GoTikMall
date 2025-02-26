package service

import (
	"context"

	"github.com/XJTU-zxc/GoTikMall/app/product/biz/dal/mysql"
	"github.com/XJTU-zxc/GoTikMall/app/product/biz/model"
	product "github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/product"
)

type SearchProductsService struct {
	ctx context.Context
} // NewSearchProductsService new SearchProductsService
func NewSearchProductsService(ctx context.Context) *SearchProductsService {
	return &SearchProductsService{ctx: ctx}
}

// Run create note info
func (s *SearchProductsService) Run(req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	// Finish your business logic.
	productQuery := model.NewProductQuery(s.ctx, mysql.DB)
	products, err := productQuery.SearchProducts(req.Query)
	if err != nil {
		return nil, err
	}
	resp = &product.SearchProductsResp{}
	for _, p := range products {
		resp.Results = append(resp.Results, &product.Product{
			Id:          uint32(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Picture:     p.Picture,
		})
	}
	return resp, nil
}
