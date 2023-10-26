package productMembershipRepository

import (
	"github.com/raj847/togrpc/models"
	"github.com/raj847/togrpc/repositories"
)

type productMembershipRepository struct {
	RepoDB repositories.Repository
}

// NewProductRepository
func NewProductMembershipRepository(repoDB repositories.Repository) productMembershipRepository {
	return productMembershipRepository{
		RepoDB: repoDB,
	}
}

func (ctx productMembershipRepository) FindProductMembershipById(id int64) (*models.ProductMembership, error) {
	var result models.ProductMembership

	query := `SELECT id, product_membership_code, product_membership_name, ou_id, product_id, 
       			due_date, disc_type, disc_amount, disc_pct, grace_period_date,
       			price, active, service_fee, is_pct_sfee, create_username, 
       			created_at, update_username, updated_at 
			FROM m_product_membership
			WHERE id = $1`

	err := ctx.RepoDB.DB.QueryRowContext(ctx.RepoDB.Context, query, id).
		Scan(&result.ID, &result.ProductMembershipCode, &result.ProductMembershipName, &result.OuId, &result.ProductId,
			&result.DueDate, &result.DiscType, &result.DiscAmount, &result.DiscPct, &result.GracePeriodDate, &result.Price,
			&result.Active, &result.ServiceFee, &result.IsPct, &result.CreateUsername,
			&result.CreatedAt, &result.UpdateUsername, &result.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
