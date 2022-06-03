package gateway

import (
	"context"
	"minimarket/pkg/rest"
	"time"
)

type Service interface {
	Register(context.Context, string, string) error
	Authorize(context.Context, string, string) (*Auth, error)
	GetProducts(context.Context, rest.QueryParameter, string, string) ([]Product, error)
	WriteComment(context.Context, int, int, string, string) error
}

type GatewayService struct {
	UserCl    UserClientInterface
	ProductCl ProductClientInterface
}

func New(usercl UserClientInterface, productcl ProductClientInterface) Service {
	return &GatewayService{
		UserCl:    usercl,
		ProductCl: productcl,
	}
}

func (svc *GatewayService) Register(ctx context.Context, email string, password string) error {
	_, err := svc.UserCl.RegisterUser(ctx, email, password)
	if err != nil {
		return err
	}

	return nil
}

func (svc *GatewayService) Authorize(ctx context.Context, email string, password string) (*Auth, error) {
	auth, err := svc.UserCl.AuthorizeUser(ctx, email, password)
	if err != nil {
		return nil, err
	}
	return &Auth{
		Token: auth.Token,
	}, nil
}

func (svc *GatewayService) GetProducts(ctx context.Context, query rest.QueryParameter, name string, category string) ([]Product, error) {
	var products []Product

	res, err := svc.ProductCl.GetProducts(ctx, name, category, query)
	if err != nil {
		return nil, err
	}

	replyMap := map[int32][]Reply{}
	for _, r := range res.Comments {
		if r.Replyid != 0 {
			if _, ok := replyMap[r.Replyid]; !ok {
				replyMap[r.Replyid] = []Reply{}
			}

			t, err := time.Parse(time.RFC3339, r.Createdat)
			if err != nil {
				return nil, err
			}

			replyMap[r.Replyid] = append(replyMap[r.Replyid], Reply{
				ID:        int(r.Id),
				ProductID: int(r.Productid),
				Message:   r.Message,
				Owner:     r.Owner,
				CreatedAt: t,
			})
		}
	}

	commentMap := map[int32][]ProductComment{}
	for _, c := range res.Comments {
		if c.Replyid == 0 {
			if _, ok := commentMap[c.Productid]; !ok {
				commentMap[c.Productid] = []ProductComment{}
			}

			t, err := time.Parse(time.RFC3339, c.Createdat)
			if err != nil {
				return nil, err
			}

			commentMap[c.Productid] = append(commentMap[c.Productid], ProductComment{
				ID:        int(c.Id),
				ProductID: int(c.Productid),
				Message:   c.Message,
				Owner:     c.Owner,
				Replies:   replyMap[c.Id],
				CreatedAt: t,
			})
		}
	}

	for _, p := range res.Products {
		product := Product{
			ID:                 int(p.Id),
			ProductName:        p.Name,
			ProductDescription: p.Description,
			ProductCategory:    p.Category,
			ProductPhoto:       p.Picture,
			ProductComment:     commentMap[p.Id],
		}

		products = append(products, product)
	}

	return products, nil
}

func (svc *GatewayService) WriteComment(ctx context.Context, productID int, ReplyID int, message string, owner string) error {
	_, err := svc.ProductCl.WriteComment(ctx, productID, ReplyID, owner, message)
	if err != nil {
		return err
	}

	return nil
}
