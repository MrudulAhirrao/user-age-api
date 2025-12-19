package service

import (
    "context"
    "errors"
    "time"

    "github.com/jackc/pgx/v5/pgtype"
    db "user-age-api/db/sqlc" 
)

// The Request format (JSON)
type UpdateProfileRequest struct {
    UserID  int32  `json:"-"`
    Name    string `json:"name"`
    Dob     string `json:"dob"`
    Address struct {
        Line1      string `json:"line1"`
        Line2      string `json:"line2"`
        City       string `json:"city"`
        State      string `json:"state"`
        PostalCode string `json:"postal_code"`
        Country    string `json:"country"`
    } `json:"address"`
}

func (s *AuthService) UpdateProfile(ctx context.Context, req UpdateProfileRequest) error {
    // 1. Start Transaction
    tx, err := s.db.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    qTx := s.queries.WithTx(tx)

    // 2. Prepare the Date (Logic to handle empty vs new date)
    var dobParams pgtype.Date
    if req.Dob != "" {
        // Parse the string "1990-01-01" into a Time object
        parsedTime, err := time.Parse("2006-01-02", req.Dob)
        if err != nil {
            return errors.New("invalid dob format, use YYYY-MM-DD")
        }
        // Valid Date -> Send it to DB
        dobParams = pgtype.Date{Time: parsedTime, Valid: true}
    } else {
        // Empty String -> Send NULL to DB (SQL will keep old value)
        dobParams = pgtype.Date{Valid: false}
    }

    // 3. Update User (Name AND Dob)
    _, err = qTx.UpdateUser(ctx, db.UpdateUserParams{
        ID:   req.UserID,
        Name: req.Name,
        Dob:  dobParams,
    })
    if err != nil {
        return err
    }

    // 4. Update Address (No changes needed here if it was working)
    _, err = qTx.CreateOrUpdateAddress(ctx, db.CreateOrUpdateAddressParams{
        UserID:     req.UserID,
        Line1:      req.Address.Line1,
        Line2:      pgtype.Text{String: req.Address.Line2, Valid: req.Address.Line2 != ""},
        City:       req.Address.City,
        State:      req.Address.State,
        PostalCode: req.Address.PostalCode,
        Country:    req.Address.Country,
    })
    if err != nil {
        return err
    }

    return tx.Commit(ctx)
}