package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/apierror"
	"github.com/dapperAuteur/dashboard-go-api/internal/budget"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
)

// Vendor defines all of the handlers related to vendor.
// It holds the application state needed by the handler methods.
type Vendor struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// ListVendors gets all vendors from the service layer.
func (v Vendor) ListVendors(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Vendor.ListVendors")
	defer span.End()

	list, err := budget.ListVendors(ctx, v.DB)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

// CreateVendor decodes the body of a request to create a new vendor.
// The full vendor with generated fields is sent back in the response.
func (v Vendor) CreateVendor(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("auth claims missing from context")
	}

	var newVendor budget.NewVendor

	if err := web.Decode(r, &newVendor); err != nil {
		return err
	}

	vendorCreated, err := budget.CreateVendor(ctx, v.DB, claims, newVendor, time.Now())
	if err != nil {
		switch err {
		case apierror.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "creating vendor %q", newVendor)
		}
	}
	return web.Respond(ctx, w, vendorCreated, http.StatusCreated)
}

// RetrieveVendor will get the vendor from the db identified by an _id in the request URL, then encodes it in a response client.
func (v Vendor) RetrieveVendor(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_id := chi.URLParam(r, "_id")

	vFound, err := budget.RetrieveVendor(ctx, v.DB, _id)
	if err != nil {
		switch err {
		case apierror.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case apierror.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for vendor %q", _id)
		}
	}
	return web.Respond(ctx, w, vFound, http.StatusOK)
}

// UpdateOneVendor decodes the body of a request to update an existing financial account.
// The _id of the financial account is part of the request URL.
func (v *Vendor) UpdateOneVendor(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	vendorID := chi.URLParam(r, "_id")

	var vendorUpdate budget.UpdateVendor
	if err := web.Decode(r, &vendorUpdate); err != nil {
		return errors.Wrap(err, "decoding vendor update")
	}

	if err := budget.UpdateOneVendor(ctx, v.DB, claims, vendorID, vendorUpdate, time.Now()); err != nil {
		switch err {
		case apierror.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case apierror.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case apierror.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "updating vendor %q", vendorID)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusOK)
}

// DeleteVendor removes a single vendor identified by a vendor ID inthe request URL
func (v *Vendor) DeleteVendor(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	vendorID := chi.URLParam(r, "_id")

	if err := budget.DeleteVendor(ctx, v.DB, claims, vendorID); err != nil {
		switch err {
		case apierror.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case apierror.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case apierror.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "deleting vendor %q", vendorID)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
