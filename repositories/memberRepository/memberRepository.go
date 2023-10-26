package memberRepository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"github.com/raj847/togrpc/constans"
	"github.com/raj847/togrpc/models"
	"github.com/raj847/togrpc/repositories"
	"github.com/raj847/togrpc/utils"
)

type memberRepository struct {
	RepoDB repositories.Repository
}

func NewMemberRepository(repo repositories.Repository) memberRepository {
	return memberRepository{
		RepoDB: repo,
	}
}

func (ctx memberRepository) AddMember(requestMember models.Member, tx *sql.Tx) (int64, error) {
	var err error
	var id int64

	query := `INSERT INTO member (
		partner_code, first_name, last_name, role_type, phone_number,
		 email, active, active_at, non_active_at, ou_id, 
		 type_partner, card_number, vehicle_number, registered_datetime, date_from,
		 date_to, product_id, product_code, created_at, created_by, 
		 updated_at, updated_by
	) VALUES (
		$1, $2, $3, $4, $5,
		$6, $7, $8, $9, $10,
		$11, $12, $13, $14, $15,
		$16, $17, $18, $19, $20,
		$21, $22
	) RETURNING id`

	if tx != nil {
		err = tx.QueryRow(query, requestMember.PartnerCode, requestMember.FirstName, requestMember.LastName,
			requestMember.RoleType, requestMember.PhoneNumber, requestMember.Email, requestMember.Active,
			requestMember.ActiveAt, requestMember.NonActiveAt, requestMember.OuId, requestMember.TypePartner,
			requestMember.CardNumber, requestMember.VehicleNumber, requestMember.RegisteredDatetime,
			requestMember.DateFrom, requestMember.DateTo, requestMember.ProductId, requestMember.ProductCode,
			requestMember.CreatedAt, requestMember.CreatedBy, requestMember.UpdatedAt, requestMember.UpdatedBy).Scan(&id)
	} else {
		err = ctx.RepoDB.DB.QueryRow(query, requestMember.PartnerCode, requestMember.FirstName, requestMember.LastName,
			requestMember.RoleType, requestMember.PhoneNumber, requestMember.Email, requestMember.Active,
			requestMember.ActiveAt, requestMember.NonActiveAt, requestMember.OuId, requestMember.TypePartner,
			requestMember.CardNumber, requestMember.VehicleNumber, requestMember.RegisteredDatetime,
			requestMember.DateFrom, requestMember.DateTo, requestMember.ProductId, requestMember.ProductCode,
			requestMember.CreatedAt, requestMember.CreatedBy, requestMember.UpdatedAt, requestMember.UpdatedBy).Scan(&id)
	}

	if err != nil {
		return id, err
	}

	return id, nil
}

func (ctx memberRepository) IsPartnerExistsByCardNumber(cardNumber string) (models.ResponseFindPartner, bool) {
	var result models.ResponseFindPartner

	var query = `
	SELECT id, date_from, date_to, card_number, ou_id, product_id, partner_code
	FROM member 
	WHERE card_number = $1
    GROUP BY id
	ORDER BY MAX(id) DESC
	LIMIT 1`

	rows, err := ctx.RepoDB.DB.Query(query, cardNumber)
	if err != nil {
		return result, false
	}
	defer rows.Close()

	data, _ := memberSpecialDto(rows)
	if len(data) == 0 {
		return result, false
	}

	return data[0], true

}

func (ctx memberRepository) IsPartnerExistsById(id int64) (models.ResponseIsPartnerExistsByID, bool) {
	var count int64
	var result models.ResponseIsPartnerExistsByID

	var query = `
	SELECT partner_code, date_from, date_to, active_at, COUNT(1)
	FROM member 
	WHERE id = $1
    GROUP BY id
	ORDER BY MAX(id) DESC
	LIMIT 1`

	err := ctx.RepoDB.DB.QueryRow(query, id).Scan(&result.PartnerCode, &result.DateFrom, &result.DateTo, &result.ActiveAt, &count)
	if err != nil {
		return result, false
	}

	if count > 0 {
		return result, false
	}

	return result, true

}

func (ctx memberRepository) IsMemberActiveExistsByUUIDCard(cardNumber string, date string) (models.Member, bool) {
	var result models.Member

	var query = `
	SELECT id, partner_code, first_name, last_name, role_type, phone_number,
		email, active, active_at, non_active_at, ou_id, 
		type_partner, card_number, vehicle_number, registered_datetime, date_from,
		date_to, product_id, product_code, created_at, created_by, 
		updated_at, updated_by
	FROM member
	WHERE card_number = $1 
	  AND $2 BETWEEN date_from AND date_to
	  AND active = $3`

	rows, err := ctx.RepoDB.DB.Query(query, cardNumber, date, constans.YES)
	if err != nil {
		return result, false
	}
	defer rows.Close()

	data, _ := memberDto(rows)
	if len(data) == 0 {
		return result, false
	}

	return data[0], true
}

func (ctx memberRepository) IsCardNumberExists(CardNumber string) bool {
	var count int64

	var query = `SELECT COUNT(1) 
	FROM member
	WHERE card_number = $1`

	err := ctx.RepoDB.DB.QueryRow(query, CardNumber).Scan(&count)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if count > 0 {
		return false
	}

	return true
}

func (ctx memberRepository) ActivationPartner(updatePartner models.EditPartner, tx *sql.Tx) error {
	var err error

	query := `
		UPDATE member SET active = $1 , active_at = $2, non_active_at = $3
		WHERE id = $4 `

	if tx != nil {
		_, err = tx.Query(query, updatePartner.Active, updatePartner.ActiveAt, updatePartner.NonActiveAt, updatePartner.ID)
		if err != nil {
			return err
		}

	} else {
		_, err = ctx.RepoDB.DB.Query(query, updatePartner.Active, updatePartner.ActiveAt, updatePartner.NonActiveAt, updatePartner.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ctx memberRepository) GetListPartnerAdvance(requestPartner models.RequestFindPartnerAdvance) ([]models.Member, error) {
	var result []models.Member
	var args []interface{}

	var query = `
	SELECT id, partner_code, first_name, last_name, role_type, phone_number,
	email, active, active_at, non_active_at, ou_id, 
	type_partner, card_number, vehicle_number, registered_datetime, date_from,
	date_to, product_id, product_code, created_at, created_by, 
	updated_at, updated_by
	FROM member
	WHERE true`

	if requestPartner.Keyword != constans.EMPTY_VALUE {
		query += ` AND (partner_code ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?) `
		args = append(args, "%"+requestPartner.Keyword+"%", "%"+requestPartner.Keyword+"%", "%"+requestPartner.Keyword+"%")
	}

	if requestPartner.ColumnOrderName != constans.EMPTY_VALUE {
		if requestPartner.AscDesc == constans.ASCENDING {
			query += ` ORDER BY ` + requestPartner.ColumnOrderName + ` ASC `
		} else if requestPartner.AscDesc == constans.DESCENDING {
			query += ` ORDER BY ` + requestPartner.ColumnOrderName + ` DESC `
		} else {
			query += ` ORDER BY partner_code DESC`
		}
	}

	query += ` LIMIT ? OFFSET ? `
	args = append(args, requestPartner.Limit, requestPartner.Offset)

	newQuery := utils.ReplaceSQL(query, "?")
	rows, err := ctx.RepoDB.DB.Query(newQuery, args...)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		return result, err
	}

	return memberDto(rows)
}

func (ctx memberRepository) CountGetListPartnerAdvance(requestPartner models.RequestFindPartnerAdvance) (int64, error) {
	var count int64
	var args []interface{}

	var query = `
	SELECT COUNT(1)
	FROM member
	WHERE true`

	if requestPartner.Keyword != constans.EMPTY_VALUE {
		query += ` AND (partner_code ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?) `
		args = append(args, "%"+requestPartner.Keyword+"%", "%"+requestPartner.Keyword+"%", "%"+requestPartner.Keyword+"%")
	}

	if requestPartner.ColumnOrderName != constans.EMPTY_VALUE {
		if requestPartner.AscDesc == constans.ASCENDING {
			query += ` ORDER BY ` + requestPartner.ColumnOrderName + ` ASC `
		} else if requestPartner.AscDesc == constans.DESCENDING {
			query += ` ORDER BY ` + requestPartner.ColumnOrderName + ` DESC `
		} else {
			query += ` ORDER BY partner_code DESC`
		}
	}

	query += ` LIMIT ? OFFSET ? `
	args = append(args, requestPartner.Limit, requestPartner.Offset)

	newQuery := utils.ReplaceSQL(query, "?")
	err := ctx.RepoDB.DB.QueryRow(newQuery, args...).Scan(&count)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (ctx memberRepository) IsMemberFreePassByIndex(uuidCard string) (*models.Member, bool, error) {
	var result []models.Member

	query := `
		SELECT partner_code, first_name, last_name, role_type, phone_number,
			email, active, active_at, non_active_at, ou_id, 
			type_partner, card_number, vehicle_number, registered_datetime, date_from, 
			date_to, product_id, product_code, created_at, created_by, 
			updated_at, updated_by
		FROM member
		WHERE card_number = $1
			AND active = $2
			AND type_partner = $3`

	rows, err := ctx.RepoDB.DB.Query(query, uuidCard, constans.YES, constans.TYPE_PARTNER_FREE_PASS)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	for rows.Next() {
		var val models.Member
		err = rows.Scan(&val.PartnerCode, &val.FirstName, &val.LastName, &val.RoleType, &val.PhoneNumber,
			&val.Email, &val.Active, &val.ActiveAt, &val.NonActiveAt, &val.OuId,
			&val.TypePartner, &val.CardNumber, &val.VehicleNumber, &val.RegisteredDatetime, &val.DateFrom,
			&val.DateTo, &val.ProductId, &val.ProductCode, &val.CreatedAt, &val.CreatedBy,
			&val.UpdatedAt, &val.UpdatedBy)

		if err != nil {
			return nil, false, err
		}

		result = append(result, val)
	}

	if len(result) == 0 {
		return nil, false, nil
	}

	return &result[0], true, nil
}

func (ctx memberRepository) IsMemberByAdvanceIndex(uuidCard, vehicleNumber, inquiryDate string, memberBy string, isFreePass bool) (*models.Member, bool, error) {
	var result []models.Member
	var args []interface{}

	query := `
		SELECT partner_code, first_name, last_name, role_type, phone_number,
			email, active, active_at, non_active_at, ou_id, 
			type_partner, card_number, vehicle_number, registered_datetime, date_from, 
			date_to, product_id, product_code, created_at, created_by, 
			updated_at, updated_by
		FROM member 
		WHERE true 
			AND active = ?`

	args = append(args, constans.YES)

	if isFreePass {
		query += ` AND type_partner = ? `
		args = append(args, constans.TYPE_PARTNER_FREE_PASS)
	} else {
		query += ` AND ? BETWEEN date_from AND date_to `
		args = append(args, inquiryDate)
	}

	if memberBy == constans.VALIDATE_MEMBER_CARD {
		query += ` AND card_number = ? `
		args = append(args, uuidCard)
	}

	if memberBy == constans.VALIDATE_MEMBER_NOPOL {
		query += ` AND vehicle_number = ? `
		args = append(args, vehicleNumber)
	}

	if memberBy == constans.VALIDATE_MEMBER_MIX {
		query += ` AND vehicle_number = ? 
						AND card_number = ? `
		args = append(args, vehicleNumber, uuidCard)
	}

	newQuery := utils.ReplaceSQL(query, "?")

	log.Println("newQuery:", newQuery)

	rows, err := ctx.RepoDB.DB.QueryContext(ctx.RepoDB.Context, newQuery, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	for rows.Next() {
		var val models.Member
		err = rows.Scan(&val.PartnerCode, &val.FirstName, &val.LastName, &val.RoleType, &val.PhoneNumber,
			&val.Email, &val.Active, &val.ActiveAt, &val.NonActiveAt, &val.OuId,
			&val.TypePartner, &val.CardNumber, &val.VehicleNumber, &val.RegisteredDatetime, &val.DateFrom,
			&val.DateTo, &val.ProductId, &val.ProductCode, &val.CreatedAt, &val.CreatedBy,
			&val.UpdatedAt, &val.UpdatedBy)

		if err != nil {
			return nil, false, err
		}

		result = append(result, val)
	}

	if len(result) == 0 {
		return nil, false, nil
	}

	return &result[0], true, nil
}

func (ctx memberRepository) IsMemberByAdvanceIndexCustom(uuidCard, vehicleNumber, inquiryDate string, memberBy string, isFreePass bool) (*models.Member, bool, error) {
	var result []models.Member
	var args []interface{}

	query := `
		SELECT partner_code, first_name, last_name, role_type, phone_number,
			email, active, active_at, non_active_at, ou_id, 
			type_partner, card_number, vehicle_number, registered_datetime, date_from, 
			date_to, product_id, product_code, created_at, created_by, 
			updated_at, updated_by
		FROM member 
		WHERE true 
			AND active = ?`

	args = append(args, constans.YES)

	log.Println(inquiryDate)
	if isFreePass {
		query += ` AND type_partner = ? `
		args = append(args, constans.TYPE_PARTNER_FREE_PASS)
	} else {
		query += ` AND ? BETWEEN LEFT(date_from,10) AND LEFT(date_to,10) `
		args = append(args, inquiryDate)
	}

	if memberBy == constans.VALIDATE_MEMBER_CARD {
		query += ` AND card_number = ? `
		args = append(args, uuidCard)
	}

	if memberBy == constans.VALIDATE_MEMBER_NOPOL {
		query += ` AND vehicle_number = ? `
		args = append(args, vehicleNumber)
	}

	if memberBy == constans.VALIDATE_MEMBER_MIX {
		query += ` AND vehicle_number = ? 
						AND card_number = ? `
		args = append(args, vehicleNumber, uuidCard)
	}

	newQuery := utils.ReplaceSQL(query, "?")

	log.Println("newQuery:", newQuery)

	rows, err := ctx.RepoDB.DB.QueryContext(ctx.RepoDB.Context, newQuery, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	for rows.Next() {
		var val models.Member
		err = rows.Scan(&val.PartnerCode, &val.FirstName, &val.LastName, &val.RoleType, &val.PhoneNumber,
			&val.Email, &val.Active, &val.ActiveAt, &val.NonActiveAt, &val.OuId,
			&val.TypePartner, &val.CardNumber, &val.VehicleNumber, &val.RegisteredDatetime, &val.DateFrom,
			&val.DateTo, &val.ProductId, &val.ProductCode, &val.CreatedAt, &val.CreatedBy,
			&val.UpdatedAt, &val.UpdatedBy)

		if err != nil {
			return nil, false, err
		}

		result = append(result, val)
	}

	if len(result) == 0 {
		return nil, false, nil
	}

	return &result[0], true, nil
}

func (ctx memberRepository) GetMemberActiveListByPeriod(uuidCard, vehicleNumber, checkinDate string, inquiryDate string, memberBy string, isFreePass bool) ([]models.Member, error) {
	var result []models.Member
	var args []interface{}

	query := `
		SELECT partner_code, first_name, last_name, role_type, phone_number,
			email, active, active_at, non_active_at, ou_id, 
			type_partner, card_number, vehicle_number, registered_datetime, date_from, 
			date_to, product_id, product_code, created_at, created_by, 
			updated_at, updated_by
		FROM member 
		WHERE true 
			AND active = ?`

	args = append(args, constans.YES)

	if isFreePass {
		query += ` AND type_partner = ? `
		args = append(args, constans.TYPE_PARTNER_FREE_PASS)
	} else {
		query += ` AND ( ? BETWEEN date_from AND date_to OR ? BETWEEN date_from AND date_to ) `
		args = append(args, checkinDate, inquiryDate)
	}

	if memberBy == constans.VALIDATE_MEMBER_CARD {
		query += ` AND card_number = ? `
		args = append(args, uuidCard)
	}

	if memberBy == constans.VALIDATE_MEMBER_NOPOL {
		query += ` AND vehicle_number = ? `
		args = append(args, vehicleNumber)
	}

	if memberBy == constans.VALIDATE_MEMBER_MIX {
		query += ` AND vehicle_number = ? 
						AND card_number = ? `
		args = append(args, vehicleNumber, uuidCard)
	}

	query += ` ORDER BY date_to DESC `

	newQuery := utils.ReplaceSQL(query, "?")
	rows, err := ctx.RepoDB.DB.QueryContext(ctx.RepoDB.Context, newQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var val models.Member
		err = rows.Scan(&val.PartnerCode, &val.FirstName, &val.LastName, &val.RoleType, &val.PhoneNumber,
			&val.Email, &val.Active, &val.ActiveAt, &val.NonActiveAt, &val.OuId,
			&val.TypePartner, &val.CardNumber, &val.VehicleNumber, &val.RegisteredDatetime, &val.DateFrom,
			&val.DateTo, &val.ProductId, &val.ProductCode, &val.CreatedAt, &val.CreatedBy,
			&val.UpdatedAt, &val.UpdatedBy)

		if err != nil {
			return nil, err
		}

		result = append(result, val)
	}

	return result, nil
}

func (ctx memberRepository) GetMemberActiveListByPeriodCustom(uuidCard, vehicleNumber, checkinDate string, inquiryDate string, memberBy string, isFreePass bool, isHotelMember bool) ([]models.Member, error) {
	var result []models.Member
	var args []interface{}

	query := `
		SELECT partner_code, first_name, last_name, role_type, phone_number,
			email, active, active_at, non_active_at, ou_id, 
			type_partner, card_number, vehicle_number, registered_datetime, date_from, 
			date_to, product_id, product_code, created_at, created_by, 
			updated_at, updated_by
		FROM member 
		WHERE true 
			AND active = ?`

	args = append(args, constans.YES)

	if isFreePass {
		query += ` AND type_partner = ? `
		args = append(args, constans.TYPE_PARTNER_FREE_PASS)
	} else if isHotelMember {
		query += ` AND ( ? BETWEEN LEFT(date_from,16) AND LEFT(date_to,16) OR ? BETWEEN LEFT(date_from,16) AND LEFT(date_to,16) ) `
		args = append(args, checkinDate, inquiryDate)
	} else {
		query += ` AND ( ? BETWEEN LEFT(date_from,10) AND LEFT(date_to,10) OR ? BETWEEN LEFT(date_from,10) AND LEFT(date_to,10) ) `
		args = append(args, checkinDate, inquiryDate)
	}

	if memberBy == constans.VALIDATE_MEMBER_CARD {
		query += ` AND card_number = ? `
		args = append(args, uuidCard)
	}

	if memberBy == constans.VALIDATE_MEMBER_NOPOL {
		query += ` AND vehicle_number = ? `
		args = append(args, vehicleNumber)
	}

	if memberBy == constans.VALIDATE_MEMBER_MIX {
		query += ` AND vehicle_number = ? 
						AND card_number = ? `
		args = append(args, vehicleNumber, uuidCard)
	}

	query += ` ORDER BY date_to DESC `

	newQuery := utils.ReplaceSQL(query, "?")
	rows, err := ctx.RepoDB.DB.QueryContext(ctx.RepoDB.Context, newQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var val models.Member
		err = rows.Scan(&val.PartnerCode, &val.FirstName, &val.LastName, &val.RoleType, &val.PhoneNumber,
			&val.Email, &val.Active, &val.ActiveAt, &val.NonActiveAt, &val.OuId,
			&val.TypePartner, &val.CardNumber, &val.VehicleNumber, &val.RegisteredDatetime, &val.DateFrom,
			&val.DateTo, &val.ProductId, &val.ProductCode, &val.CreatedAt, &val.CreatedBy,
			&val.UpdatedAt, &val.UpdatedBy)

		if err != nil {
			return nil, err
		}

		result = append(result, val)
	}

	return result, nil
}

func (ctx memberRepository) RemoveMember(id int64, tx *sql.Tx) error {
	var err error

	query := `DELETE FROM member 
	WHERE id = $1`

	if tx != nil {
		_, err = tx.Query(query, id)
		if err != nil {
			return err
		}
	} else {
		_, err = ctx.RepoDB.DB.Query(query, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ctx memberRepository) FindMemberActiveByDate(Date string, RoleType string) (models.SummaryMember, error) {
	var result models.SummaryMember

	query := `SELECT COUNT(1)
	FROM member
	WHERE active = $1
	AND $2 BETWEEN date_from AND date_to
	AND role_type = $3
	GROUP BY role_type`

	err := ctx.RepoDB.DB.QueryRow(query, constans.YES, Date, RoleType).Scan(&result.Member)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (ctx memberRepository) IsMemberExistsByPartnerCode(partnerCode, startDate, keyword string) (*models.Member, bool, error) {
	var args []interface{}
	var result models.Member

	query := `
		SELECT partner_code, first_name, last_name, role_type, phone_number,
			email, active, active_at, non_active_at, ou_id, 
			type_partner, card_number, vehicle_number, registered_datetime, date_from, 
			date_to, product_id, product_code, created_at, created_by, 
			updated_at, updated_by
		FROM member 
		WHERE partner_code = ?
			AND ? BETWEEN date_from AND date_to 
			AND card_number = ?
			AND type_partner = ?`
	args = append(args, partnerCode, startDate, keyword, constans.TYPE_PARTNER_ONE_TIME)

	newQuery := utils.ReplaceSQL(query, "?")
	err := ctx.RepoDB.DB.QueryRow(newQuery, args...).Scan(&result.PartnerCode, &result.FirstName, &result.LastName, &result.RoleType, &result.PhoneNumber,
		&result.Email, &result.Active, &result.ActiveAt, &result.NonActiveAt, &result.OuId,
		&result.TypePartner, &result.CardNumber, &result.VehicleNumber, &result.RegisteredDatetime, &result.DateFrom,
		&result.DateTo, &result.ProductId, &result.ProductCode, &result.CreatedAt, &result.CreatedBy,
		&result.UpdatedAt, &result.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		} else {
			return nil, false, err
		}
	}

	return &result, true, nil
}

func (ctx memberRepository) IsMemberActiveByCustom(uuidCard, vehicleNumber, inquiryDate string, memberBy string, isFreePass bool) (*models.MemberCustom, bool, error) {
	var result []models.MemberCustom
	var args []interface{}

	query := `
		SELECT partner_code, first_name, last_name, role_type, phone_number,
			CASE 
				WHEN ? BETWEEN date_from AND date_to  THEN 'ACTIVE'
				WHEN date_from > ? AND date_to > ? THEN 'PRE ACTIVE'
				ELSE 'EXPIRED'
			END as status,
			email, active, active_at, non_active_at, ou_id, 
			type_partner, card_number, vehicle_number, registered_datetime, date_from, 
			date_to, product_id, product_code, created_at, created_by, 
			updated_at, updated_by
		FROM member 
		WHERE true 
			AND active = ?`

	args = append(args, inquiryDate, inquiryDate, inquiryDate, constans.YES)

	if isFreePass {
		query += ` AND type_partner = ? `
		args = append(args, constans.TYPE_PARTNER_FREE_PASS)
	} else {
		query += ` AND ? BETWEEN date_from AND date_to `
		args = append(args, inquiryDate)
	}

	if memberBy == constans.VALIDATE_MEMBER_CARD {
		query += ` AND card_number = ? `
		args = append(args, uuidCard)
	}

	if memberBy == constans.VALIDATE_MEMBER_NOPOL {
		query += ` AND vehicle_number = ? `
		args = append(args, vehicleNumber)
	}

	if memberBy == constans.VALIDATE_MEMBER_MIX {
		query += ` AND vehicle_number = ? 
						AND card_number = ? `
		args = append(args, vehicleNumber, uuidCard)
	}

	query += ` ORDER BY date_to DESC `

	newQuery := utils.ReplaceSQL(query, "?")

	log.Println("newQuery:", newQuery)

	rows, err := ctx.RepoDB.DB.QueryContext(ctx.RepoDB.Context, newQuery, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	for rows.Next() {
		var val models.MemberCustom
		err = rows.Scan(&val.PartnerCode, &val.FirstName, &val.LastName, &val.RoleType, &val.PhoneNumber,
			&val.Status, &val.Email, &val.Active, &val.ActiveAt, &val.NonActiveAt, &val.OuId,
			&val.TypePartner, &val.CardNumber, &val.VehicleNumber, &val.RegisteredDatetime, &val.DateFrom,
			&val.DateTo, &val.ProductId, &val.ProductCode, &val.CreatedAt, &val.CreatedBy,
			&val.UpdatedAt, &val.UpdatedBy)

		if err != nil {
			return nil, false, err
		}

		result = append(result, val)

	}

	if len(result) == 0 {
		return nil, false, nil
	}

	return &result[0], true, nil
}

func (ctx memberRepository) GetListMemberActiveByCustom(uuidCard, vehicleNumber, inquiryDate string, memberBy string, isFreePass bool) (*[]models.MemberCustom, bool, error) {
	var result []models.MemberCustom
	var args []interface{}

	query := `
		SELECT partner_code, first_name, last_name, role_type, phone_number,
			CASE 
				WHEN ? BETWEEN date_from AND date_to  THEN 'ACTIVE'
				WHEN date_from > ? AND date_to > ? THEN 'PRE ACTIVE'
				ELSE 'EXPIRED'
			END as status,
			email, active, active_at, non_active_at, ou_id, 
			type_partner, card_number, vehicle_number, registered_datetime, date_from, 
			date_to, product_id, product_code, created_at, created_by, 
			updated_at, updated_by
		FROM member 
		WHERE true 
			AND active = ?`

	args = append(args, inquiryDate, inquiryDate, inquiryDate, constans.YES)

	query += ` AND ? BETWEEN date_from AND date_to `
	args = append(args, inquiryDate)

	if memberBy == constans.VALIDATE_MEMBER_CARD {
		query += ` AND card_number = ? `
		args = append(args, uuidCard)
	}

	if memberBy == constans.VALIDATE_MEMBER_NOPOL {
		query += ` AND vehicle_number = ? `
		args = append(args, vehicleNumber)
	}

	if memberBy == constans.VALIDATE_MEMBER_MIX {
		query += ` AND vehicle_number = ? 
						AND card_number = ? `
		args = append(args, vehicleNumber, uuidCard)
	}

	query += ` ORDER BY date_to DESC `

	newQuery := utils.ReplaceSQL(query, "?")

	rows, err := ctx.RepoDB.DB.QueryContext(ctx.RepoDB.Context, newQuery, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	for rows.Next() {
		var val models.MemberCustom
		err = rows.Scan(&val.PartnerCode, &val.FirstName, &val.LastName, &val.RoleType, &val.PhoneNumber,
			&val.Status, &val.Email, &val.Active, &val.ActiveAt, &val.NonActiveAt, &val.OuId,
			&val.TypePartner, &val.CardNumber, &val.VehicleNumber, &val.RegisteredDatetime, &val.DateFrom,
			&val.DateTo, &val.ProductId, &val.ProductCode, &val.CreatedAt, &val.CreatedBy,
			&val.UpdatedAt, &val.UpdatedBy)

		if err != nil {
			return nil, false, err
		}

		result = append(result, val)

	}

	if len(result) == 0 {
		return nil, false, nil
	}

	return &result, true, nil
}

func (ctx memberRepository) UpdateMemberByPartnerCode(updateMember models.UpdateMember, tx *sql.Tx) error {
	var err error

	query := ` UPDATE member 
					SET date_from = $1, date_to = $2, updated_at = $3, updated_by = $4
			  WHERE partner_code = $5 AND type_partner = $6 AND card_number = $7 `

	if tx != nil {
		_, err = tx.ExecContext(ctx.RepoDB.Context, query, updateMember.DateFrom, updateMember.DateTo, updateMember.UpdatedAt,
			updateMember.UpdatedBy, updateMember.PartnerCode, constans.TYPE_PARTNER_ONE_TIME, updateMember.CardNumber)
	} else {
		_, err = ctx.RepoDB.DB.ExecContext(ctx.RepoDB.Context, query, updateMember.DateFrom, updateMember.DateTo, updateMember.UpdatedAt,
			updateMember.UpdatedBy, updateMember.PartnerCode, constans.TYPE_PARTNER_ONE_TIME, updateMember.CardNumber)
	}
	if err != nil {
		return err
	}
	return nil
}

func memberDto(rows *sql.Rows) ([]models.Member, error) {
	var result []models.Member

	for rows.Next() {
		var val models.Member
		err := rows.Scan(&val.ID, &val.PartnerCode, &val.FirstName, &val.LastName, &val.RoleType, &val.PhoneNumber,
			&val.Email, &val.Active, &val.ActiveAt, &val.NonActiveAt, &val.OuId,
			&val.TypePartner, &val.CardNumber, &val.VehicleNumber, &val.RegisteredDatetime, &val.DateFrom,
			&val.DateTo, &val.ProductId, &val.ProductCode, &val.CreatedAt, &val.CreatedBy,
			&val.UpdatedAt, &val.UpdatedBy)
		if err != nil {
			return result, err
		}
		result = append(result, val)
	}
	return result, nil
}

func memberSpecialDto(rows *sql.Rows) ([]models.ResponseFindPartner, error) {
	var result []models.ResponseFindPartner

	for rows.Next() {
		var val models.ResponseFindPartner
		err := rows.Scan(&val.PartnerId, &val.DateFrom, &val.DateTo, &val.CardNumber, &val.OuId,
			&val.ProductId, &val.PartnerCode)
		if err != nil {
			return result, err
		}
		result = append(result, val)
	}
	return result, nil
}

func (ctx memberRepository) GetPolicyOuPartnerByIndex(partner models.RequestFindPartnerAdvance) (*[]models.FindPolicyOuPartnerByIndex, error) {
	var result []models.FindPolicyOuPartnerByIndex
	var args []interface{}

	query := `
	SELECT 
		id, partner_code, first_name, last_name, ou_id, 
		product_id, registered_datetime, date_from, date_to, 
		CASE 
			WHEN ? BETWEEN date_from AND date_to  THEN 'ACTIVE'
			WHEN date_from > ? AND date_to > ? THEN 'PRE ACTIVE'
			ELSE 'EXPIRED'
		END as status, 
		type_partner, role_type,  email, phone_number, vehicle_number, 
		card_number, created_at, created_by, updated_at, updated_by, ref_partner_id
	FROM member
	WHERE ou_id = ANY(?::int[]) AND true `
	args = append(args, utils.DateNow(), utils.DateNow(), utils.DateNow(), fmt.Sprintf("%s%s%s", "{", partner.OuList, "}"))

	if partner.Keyword != constans.EMPTY_VALUE {
		query += ` AND (first_name ILIKE ? ) `
		args = append(args, "%"+partner.Keyword+"%")
	}

	if partner.ColumnOrderName != constans.EMPTY_VALUE {
		if partner.AscDesc == constans.ASCENDING {
			query += ` ORDER BY ` + partner.ColumnOrderName + ` ASC `
		} else if partner.AscDesc == constans.DESCENDING {
			query += ` ORDER BY ` + partner.ColumnOrderName + ` DESC `
		} else {
			query += ` ORDER BY  partner_code ASC`
		}
	} else {
		query += ` ORDER BY  partner_code ASC`
	}

	query += ` LIMIT ? OFFSET ? `
	args = append(args, partner.Limit, partner.Offset)

	newQuery := utils.ReplaceSQL(query, "?")
	rows, err := ctx.RepoDB.DB.QueryContext(ctx.RepoDB.Context, newQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var val models.FindPolicyOuPartnerByIndex
		err = rows.Scan(&val.ID, &val.PartnerCode, &val.FirstName, &val.LastName, &val.OuId,
			&val.ProductId, &val.RegisteredDatetime, &val.DateFrom, &val.DateTo, &val.Status,
			&val.TypePartner, &val.RoleType, &val.Email, &val.PhoneNumber, &val.VehicleNumber,
			&val.CardNumber, &val.CreatedAt, &val.CreatedBy, &val.UpdatedAt, &val.UpdatedBy, &val.RefPartnerId)
		if err != nil {
			return nil, err
		}

		result = append(result, val)
	}

	return &result, nil
}

func (ctx memberRepository) CountPolicyOuPartnerByIndex(partner models.RequestFindPartnerAdvance) (*int64, error) {
	var result int64
	var args []interface{}

	query := `
	SELECT COUNT(1) as count_member
	FROM member
	WHERE ou_id = ANY(?::int[])
	AND true `
	args = append(args, fmt.Sprintf("%s%s%s", "{", partner.OuList, "}"))

	if partner.Keyword != constans.EMPTY_VALUE {
		query += ` AND (partner_code ILIKE ? OR phone_number ILIKE ? OR first_name ILIKE ? OR vehicle_number ILIKE ? OR card_number ILIKE ?) `
		args = append(args, "%"+partner.Keyword+"%", "%"+partner.Keyword+"%", "%"+partner.Keyword+"%", "%"+partner.Keyword+"%", "%"+partner.Keyword+"%")
	}

	if partner.StatusMember == constans.ACTIVE_MEMBER {
		query += ` AND ? BETWEEN date_from AND date_to `
		args = append(args, utils.DateNow())
	}

	if partner.StatusMember == constans.PRE_ACTIVE_MEMBER {
		query += ` AND date_from > ? AND date_to > ? `
		args = append(args, utils.DateNow(), utils.DateNow())
	}

	if partner.StatusMember == constans.EXPIRED_MEMBER {
		query += ` AND date_to < ? `
		args = append(args, utils.DateNow())
	}

	newQuery := utils.ReplaceSQL(query, "?")
	err := ctx.RepoDB.DB.QueryRowContext(ctx.RepoDB.Context, newQuery, args...).Scan(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (ctx memberRepository) ValPartnerActiveByAdvance(valPartnerAdvance models.ValPartnerAdvance, productCode, typePartner string) error {
	var args []interface{}
	var count int64

	query := `
		SELECT COUNT(1)
		FROM member 
		WHERE ou_id = ?
			AND ? BETWEEN date_from AND date_to`

	args = append(args, valPartnerAdvance.OuId, valPartnerAdvance.DateNow)

	if typePartner == constans.MEMBER {
		if valPartnerAdvance.RegisteredType == constans.VEHICLE {
			query += ` AND vehicle_number = ? AND product_code = ?`
			args = append(args, valPartnerAdvance.VehicleNumber, productCode)
		} else if valPartnerAdvance.RegisteredType == constans.CARD_NUMBER {
			query += ` AND card_number = ? AND product_code = ? `
			args = append(args, valPartnerAdvance.CardNumber, productCode)
		} else if valPartnerAdvance.RegisteredType == constans.MIX {
			query += ` AND card_number = ? AND vehicle_number = ? AND product_code = ? `
			args = append(args, valPartnerAdvance.CardNumber, valPartnerAdvance.VehicleNumber, productCode)
		}
	} else if typePartner == constans.TYPE_PARTNER_FREE_PASS {
		if valPartnerAdvance.RegisteredType == constans.VEHICLE_NUMBER {
			query += ` AND vehicle_number = ? `
			args = append(args, valPartnerAdvance.VehicleNumber)
		} else if valPartnerAdvance.RegisteredType == constans.CARD_NUMBER {
			query += ` AND card_number = ? `
			args = append(args, valPartnerAdvance.CardNumber)
		} else if valPartnerAdvance.RegisteredType == constans.MIX {
			query += ` AND card_number = ? AND vehicle_number = ? `
			args = append(args, valPartnerAdvance.CardNumber, valPartnerAdvance.VehicleNumber)
		}
	}

	newQuery := utils.ReplaceSQL(query, "?")
	err := ctx.RepoDB.DB.QueryRowContext(ctx.RepoDB.Context, newQuery, args...).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			return err
		}
	}

	if count > 0 {
		return errors.New("member already active!")
	}

	return nil
}

func (ctx memberRepository) FindMemberByIndex(member models.RequestExtendMember) (*[]models.FindPolicyOuPartnerByIndex, error) {
	var result []models.FindPolicyOuPartnerByIndex
	var args []interface{}

	query := `
	SELECT 	id, partner_code, first_name, last_name, ou_id, 
		product_id, registered_datetime, date_from, date_to, 
		type_partner, role_type,  email, phone_number, vehicle_number, 
		card_number, created_at, created_by, updated_at, updated_by, ref_partner_id
	FROM member
	WHERE first_name ILIKE ? `
	args = append(args, "%"+member.MemberName+"%")

	if member.CardNumber != constans.EMPTY_VALUE {
		query += ` AND card_number = ? `
		args = append(args, member.CardNumber)
	}

	if member.VehicleNumber != constans.EMPTY_VALUE {
		query += ` AND vehicle_number = ? `
		args = append(args, member.VehicleNumber)
	}

	newQuery := utils.ReplaceSQL(query, "?")
	rows, err := ctx.RepoDB.DB.QueryContext(ctx.RepoDB.Context, newQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var val models.FindPolicyOuPartnerByIndex
		err = rows.Scan(&val.ID, &val.PartnerCode, &val.FirstName, &val.LastName, &val.OuId,
			&val.ProductId, &val.RegisteredDatetime, &val.DateFrom, &val.DateTo, &val.TypePartner,
			&val.RoleType, &val.Email, &val.PhoneNumber, &val.VehicleNumber, &val.CardNumber,
			&val.CreatedAt, &val.CreatedBy, &val.UpdatedAt, &val.UpdatedBy, &val.RefPartnerId)
		if err != nil {
			return nil, err
		}

		result = append(result, val)
	}

	return &result, nil
}

func (ctx memberRepository) EditPartnerInternal(partner models.EditPartnerInternal, tx *sql.Tx) (bool, error) {
	var err error

	query := `UPDATE member 
					SET date_from = $2, date_to = $3, role_type = $4, card_number = $5, vehicle_number = $6,
					first_name = $7, product_id = $8, product_code = $9, active_at = $10, updated_by = $11, updated_at = $12
			  WHERE id = $1`

	_, err = ctx.RepoDB.DB.QueryContext(ctx.RepoDB.Context, query, partner.ID, partner.DateFrom, partner.DateTo, partner.RoleType, partner.CardNumber,
		partner.VehicleNumber, partner.Name, partner.ProductId, partner.ProductCode, partner.ActiveAt, partner.UpdatedBy, utils.Timestamp())
	if err != nil {
		return false, err
	}

	return true, nil
}

func (ctx memberRepository) GetMemberNameList(partner models.RequestFindPartnerAdvance) (*[]models.ResponseGetMemberName, error) {
	var result []models.ResponseGetMemberName
	var args []interface{}

	query := `SELECT 
		first_name, ou_id, 
		role_type,  email, phone_number
	FROM member
	WHERE ou_id = ANY(?::int[]) AND 
	id NOT IN (
		SELECT id FROM  member
		EXCEPT SELECT MAX(id) FROM member
		GROUP BY first_name
	) AND true `

	args = append(args, fmt.Sprintf("%s%s%s", "{", partner.OuList, "}"))
	if partner.Keyword != constans.EMPTY_VALUE {
		query += ` AND ( phone_number ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ? OR role_type ILIKE ? ) `
		args = append(args, "%"+partner.Keyword+"%", "%"+partner.Keyword+"%", "%"+partner.Keyword+"%", "%"+partner.Keyword+"%")
	}

	if partner.ColumnOrderName != constans.EMPTY_VALUE {
		if partner.AscDesc == constans.ASCENDING {
			query += ` ORDER BY ` + partner.ColumnOrderName + ` ASC `
		} else if partner.AscDesc == constans.DESCENDING {
			query += ` ORDER BY ` + partner.ColumnOrderName + ` DESC `
		} else {
			query += ` ORDER BY first_name ASC`
		}
	} else {
		query += ` ORDER BY first_name ASC`
	}

	query += ` LIMIT ? OFFSET ? `
	args = append(args, partner.Limit, partner.Offset)

	newQuery := utils.ReplaceSQL(query, "?")
	rows, err := ctx.RepoDB.DB.QueryContext(ctx.RepoDB.Context, newQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var val models.ResponseGetMemberName
		err = rows.Scan(&val.FirstName, &val.OuId, &val.RoleType, &val.Email,
			&val.PhoneNumber)
		if err != nil {
			return nil, err
		}

		result = append(result, val)
	}

	return &result, nil
}

func (ctx memberRepository) EditPolicyOuPartner(policyOuPartner models.EditPolicyOuPartner, tx *sql.Tx) (bool, error) {
	var err error

	query := `UPDATE policy_ou_partner 
				SET  date_to = $2, date_from = $3
			  WHERE id = $1`

	_, err = ctx.RepoDB.DBCloud.ExecContext(ctx.RepoDB.Context, query, policyOuPartner.ID, policyOuPartner.EndDate, policyOuPartner.StartDate)
	if err != nil {
		return false, err
	}

	return true, nil
}
