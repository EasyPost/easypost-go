package easypost

import (
	"context"
	"fmt"
	"net/http"
)

// A User contains data about an EasyPost account and child accounts.
type User struct {
	ID                      string    `json:"id,omitempty" url:"id,omitempty"`
	Object                  string    `json:"object,omitempty" url:"object,omitempty"`
	ParentID                string    `json:"parent_id,omitempty" url:"parent_id,omitempty"`
	Name                    string    `json:"name,omitempty" url:"name,omitempty"`
	Email                   string    `json:"email,omitempty" url:"email,omitempty"`
	PhoneNumber             string    `json:"phone_number,omitempty" url:"phone_number,omitempty"`
	Balance                 string    `json:"balance,omitempty" url:"balance,omitempty"`
	RechargeAmount          string    `json:"recharge_amount,omitempty" url:"recharge_amount,omitempty"`
	SecondaryRechargeAmount string    `json:"secondary_recharge_amount,omitempty" url:"secondary_recharge_amount,omitempty"`
	RechargeThreshold       string    `json:"recharge_threshold,omitempty" url:"recharge_threshold,omitempty"`
	Children                []*User   `json:"children,omitempty" url:"children,omitempty"`
	APIKeys                 []*APIKey `json:"api_keys,omitempty" url:"api_keys,omitempty"`
	Verified                bool      `json:"verified,omitempty" url:"verified,omitempty"`
}

// UserOptions specifies options for creating or updating a user.
type UserOptions struct {
	ID                      string  `json:"-"` // IGNORE
	Email                   *string `json:"email,omitempty" url:"email,omitempty"`
	Password                *string `json:"password,omitempty" url:"password,omitempty"`
	PasswordConfirmation    *string `json:"password_confirmation,omitempty" url:"password_confirmation,omitempty"`
	CurrentPassword         *string `json:"current_password,omitempty" url:"current_password,omitempty"`
	Name                    *string `json:"name,omitempty" url:"name,omitempty"`
	Phone                   *string `json:"phone,omitempty" url:"phone,omitempty"`
	PhoneNumber             *string `json:"phone_number,omitempty" url:"phone_number,omitempty"`
	RechargeAmount          *string `json:"recharge_amount,omitempty" url:"recharge_amount,omitempty"`
	SecondaryRechargeAmount *string `json:"secondary_recharge_amount,omitempty" url:"secondary_recharge_amount,omitempty"`
	RechargeThreshold       *string `json:"recharge_threshold,omitempty" url:"recharge_threshold,omitempty"`
}

type userRequest struct {
	UserOptions *UserOptions `json:"user,omitempty" url:"user,omitempty"`
}

type ListChildUsersResult struct {
	Children []*User `json:"children,omitempty" url:"children,omitempty"`
	PaginatedCollection
}

// CreateUser creates a new child user.
//
//	 c := easypost.New(MyEasyPostAPIKey)
//	 opts := &easypost.UserOptions{Name: easypost.StringPtr("Child User")}
//		out, err := c.CreateUser(opts)
func (c *Client) CreateUser(in *UserOptions) (out *User, err error) {
	return c.CreateUserWithContext(context.Background(), in)
}

// CreateUserWithContext performs the same operation as CreateUser, but allows
// specifying a context that can interrupt the request.
func (c *Client) CreateUserWithContext(ctx context.Context, in *UserOptions) (out *User, err error) {
	err = c.do(ctx, http.MethodPost, "users", &userRequest{UserOptions: in}, &out)
	return
}

// GetUser retrieves a User object by ID.
func (c *Client) GetUser(userID string) (out *User, err error) {
	return c.GetUserWithContext(context.Background(), userID)
}

// GetUserWithContext performs the same operation as GetUser, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetUserWithContext(ctx context.Context, userID string) (out *User, err error) {
	err = c.do(ctx, http.MethodGet, "users/"+userID, nil, &out)
	return
}

// UpdateUser updates a user with the attributes given in the UpdateUserOptions
// parameter. If the ID field of UpdateUserOptions is empty, the operation is
// done on the current user. All other fields are updated if they are non-nil.
func (c *Client) UpdateUser(in *UserOptions) (out *User, err error) {
	return c.UpdateUserWithContext(context.Background(), in)
}

// UpdateUserWithContext performs the same operation as UpdateUser, but allows
// specifying a context that can interrupt the request.
func (c *Client) UpdateUserWithContext(ctx context.Context, in *UserOptions) (out *User, err error) {
	req := userRequest{UserOptions: in}
	path := "users"
	if in.ID != "" {
		path += "/" + in.ID
	}
	err = c.do(ctx, http.MethodPatch, path, req, &out)
	return
}

// DeleteUser removes a child user.
func (c *Client) DeleteUser(userID string) error {
	return c.DeleteUserWithContext(context.Background(), userID)
}

// DeleteUserWithContext performs the same operation as DeleteUser, but allows
// specifying a context that can interrupt the request.
func (c *Client) DeleteUserWithContext(ctx context.Context, userID string) error {
	return c.do(ctx, http.MethodDelete, "users/"+userID, nil, nil)
}

// RetrieveMe retrieves the current user.
func (c *Client) RetrieveMe() (out *User, err error) {
	return c.RetrieveMeWithContext(context.Background())
}

// RetrieveMeWithContext performs the same operation as RetrieveMe, but allows
// specifying a context that can interrupt the request.
func (c *Client) RetrieveMeWithContext(ctx context.Context) (out *User, err error) {
	err = c.do(ctx, http.MethodGet, "users", nil, &out)
	return
}

// ListChildUsers retrieves a list of child users.
func (c *Client) ListChildUsers(opts *ListOptions) (out *ListChildUsersResult, err error) {
	return c.ListChildUsersWithContext(context.Background(), opts)
}

// ListChildUsersWithContext performs the same operation as ListChildUsers, but allows
// specifying a context that can interrupt the request.
func (c *Client) ListChildUsersWithContext(ctx context.Context, opts *ListOptions) (out *ListChildUsersResult, err error) {
	err = c.do(ctx, http.MethodGet, "users/children", opts, &out)
	return
}

// GetNextChildUserPage returns the next page of child users.
func (c *Client) GetNextChildUserPage(collection *ListChildUsersResult) (out *ListChildUsersResult, err error) {
	return c.GetNextChildUserPageWithContext(context.Background(), collection)
}

// GetNextChildUserPageWithPageSize returns the next page of child users with a specific page size.
func (c *Client) GetNextChildUserPageWithPageSize(collection *ListChildUsersResult, pageSize int) (out *ListChildUsersResult, err error) {
	return c.GetNextChildUserPageWithPageSizeWithContext(context.Background(), collection, pageSize)
}

// GetNextChildUserPageWithContext performs the same operation as GetNextChildUserPage, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextChildUserPageWithContext(ctx context.Context, collection *ListChildUsersResult) (out *ListChildUsersResult, err error) {
	return c.GetNextChildUserPageWithPageSizeWithContext(ctx, collection, 0)
}

// GetNextChildUserPageWithPageSizeWithContext performs the same operation as GetNextChildUserPageWithPageSize, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextChildUserPageWithPageSizeWithContext(ctx context.Context, collection *ListChildUsersResult, pageSize int) (out *ListChildUsersResult, err error) {
	if len(collection.Children) == 0 {
		err = newEndOfPaginationError()
		return
	}
	lastID := collection.Children[len(collection.Children)-1].ID
	params, err := nextChildUserPageParameters(collection.HasMore, lastID, pageSize)
	if err != nil {
		return
	}
	return c.ListChildUsersWithContext(ctx, params)
}

// UpdateBrand updates the user brand.
func (c *Client) UpdateBrand(params map[string]interface{}, userID string) (out *Brand, err error) {
	return c.UpdateBrandWithContext(context.Background(), params, userID)
}

// UpdateBrandWithContext performs the same operation as UpdateBrand, but allows
// specifying a context that can interrupt the request.
func (c *Client) UpdateBrandWithContext(ctx context.Context, params map[string]interface{}, userID string) (out *Brand, err error) {
	newParams := map[string]interface{}{"brand": params}
	updateBrandURL := fmt.Sprintf("users/%s/brand", userID)
	err = c.do(ctx, http.MethodPatch, updateBrandURL, newParams, &out)
	return
}
