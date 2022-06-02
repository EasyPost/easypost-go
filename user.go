package easypost

import (
	"context"
	"fmt"
)

// A User contains data about an EasyPost account and child accounts.
type User struct {
	ID                      string    `json:"id,omitempty"`
	Object                  string    `json:"object,omitempty"`
	ParentID                string    `json:"parent_id,omitempty"`
	Name                    string    `json:"name,omitempty"`
	Email                   string    `json:"email,omitempty"`
	PhoneNumber             string    `json:"phone_number,omitempty"`
	Balance                 string    `json:"balance,omitempty"`
	RechargeAmount          string    `json:"recharge_amount,omitempty"`
	SecondaryRechargeAmount string    `json:"secondary_recharge_amount,omitempty"`
	RechargeThreshold       string    `json:"recharge_threshold,omitempty"`
	Children                []*User   `json:"children,omitempty"`
	APIKeys                 []*APIKey `json:"api_keys,omitempty"`
}

// UserOptions specifies options for creating or updating a user.
type UserOptions struct {
	ID                      string  `json:"-"`
	Email                   *string `json:"email,omitempty"`
	Password                *string `json:"password,omitempty"`
	PasswordConfirmation    *string `json:"password_confirmation,omitempty"`
	CurrentPassword         *string `json:"current_password,omitempty"`
	Name                    *string `json:"name,omitempty"`
	PhoneNumber             *string `json:"phone_number,omitempty"`
	RechargeAmount          *string `json:"recharge_amount,omitempty"`
	SecondaryRechargeAmount *string `json:"secondary_recharge_amount,omitempty"`
	RechargeThreshold       *string `json:"recharge_threshold,omitempty"`
}

type userRequest struct {
	UserOptions *UserOptions `json:"user,omitempty"`
}

// CreateUser creates a new child user.
//  c := easypost.New(MyEasyPostAPIKey)
//  opts := &easypost.UserOptions{Name: easypost.StringPtr("Child User")}
//	out, err := c.CreateUser(opts)
func (c *Client) CreateUser(in *UserOptions) (out *User, err error) {
	err = c.post(context.Background(), "users", &userRequest{UserOptions: in}, &out)
	return
}

// CreateUserWithContext performs the same operation as CreateUser, but allows
// specifying a context that can interrupt the request.
func (c *Client) CreateUserWithContext(ctx context.Context, in *UserOptions) (out *User, err error) {
	err = c.post(ctx, "users", &userRequest{UserOptions: in}, &out)
	return
}

// GetUser retrieves a User object by ID.
func (c *Client) GetUser(userID string) (out *User, err error) {
	err = c.get(context.Background(), "users/"+userID, &out)
	return
}

// GetUserWithContext performs the same operation as GetUser, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetUserWithContext(ctx context.Context, userID string) (out *User, err error) {
	err = c.get(ctx, "users/"+userID, &out)
	return
}

// UpdateUser updates a user with the attributes given in the UpdateUserOptions
// parameter. If the ID field of UpdateUserOptions is empty, the operation is
// done on the current user. All other fields are updated if they are non-nil.
func (c *Client) UpdateUser(in *UserOptions) (out *User, err error) {
	req := userRequest{UserOptions: in}
	path := "users"
	if in.ID != "" {
		path += "/" + in.ID
	}
	err = c.patch(context.Background(), path, req, &out)
	return
}

// UpdateUserWithContext performs the same operation as UpdateUser, but allows
// specifying a context that can interrupt the request.
func (c *Client) UpdateUserWithContext(ctx context.Context, in *UserOptions) (out *User, err error) {
	req := userRequest{UserOptions: in}
	path := "users"
	if in.ID != "" {
		path += "/" + in.ID
	}
	err = c.patch(ctx, path, req, &out)
	return
}

// DeleteUser removes a child user.
func (c *Client) DeleteUser(userID string) error {
	return c.del(context.Background(), "users/"+userID)
}

// DeleteUserWithContext performs the same operation as DeleteUser, but allows
// specifying a context that can interrupt the request.
func (c *Client) DeleteUserWithContext(ctx context.Context, userID string) error {
	return c.del(ctx, "users/"+userID)
}

// RetrieveMe retrieves the current user.
func (c *Client) RetrieveMe() (out *User, err error) {
	err = c.get(context.Background(), "users", &out)
	return
}

// RetrieveMeWithContext performs the same operation as RetrieveMe, but allows
// specifying a context that can interrupt the request.
func (c *Client) RetrieveMeWithContext(ctx context.Context) (out *User, err error) {
	err = c.get(ctx, "users", &out)
	return
}

// UpdateBrand updates the user brand.
func (c *Client) UpdateBrand(params map[string]interface{}, userID string) (out *Brand, err error) {
	newParams := map[string]interface{}{"brand": params}
	updateBrandURL := fmt.Sprintf("users/%s/brand", userID)
	err = c.patch(context.Background(), updateBrandURL, newParams, &out)
	return
}

// UpdateBrandWithContext performs the same operation as UpdateBrand, but allows
// specifying a context that can interrupt the request.
func (c *Client) UpdateBrandWithContext(ctx context.Context, params map[string]interface{}, userID string) (out *Brand, err error) {
	newParams := map[string]interface{}{"brand": params}
	updateBrandURL := fmt.Sprintf("users/%s/brand", userID)
	err = c.patch(ctx, updateBrandURL, newParams, &out)
	return
}
