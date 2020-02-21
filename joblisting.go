package quik

import (
  "context"
  "errors"

  "github.com/jinzhu/gorm"
)

type Listing struct {
  gorm.Model
  Status       string
  Expires      int
  CompanyID    string
  Technologies []Technology
  Requirements []Requirement
}

type ListingService interface {
  CreateListing(ctx context.Context, listing *Listing) error
  UpsertListing(ctx context.Context, listing *Listing) error
  GetListingsByCompanyID(ctx context.Context, id string) ([]*Listing, error)
  GetListings(ctx context.Context, queries Queries) ([]*Listing, error)
}

type Queries struct {
  Page     int
  Limit    int
  Language string
  Role     string
}

type Requirement struct {
  gorm.Model
  Description string
}

type Technology struct {
  gorm.Model
  Description
}
