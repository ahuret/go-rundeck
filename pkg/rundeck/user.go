package rundeck

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	requests "github.com/ahuret/go-rundeck/pkg/rundeck/requests"
	responses "github.com/ahuret/go-rundeck/pkg/rundeck/responses"
)

// User represents a user in rundeck
type User struct {
	Login     string
	FirstName string
	LastName  string
	Email     string
	Created   time.Time
	Updated   time.Time
	LastJob   time.Time
	Tokens    int
}

// Users represents a collection of users
type Users []User

// ListUsers returns all rundeck users
// http://rundeck.org/docs/api/index.html#list-users
func (c *Client) ListUsers() (Users, error) {
	if err := c.checkRequiredAPIVersion(responses.ListUsersResponse{}); err != nil {
		return nil, err
	}
	users := Users{}
	listusers := responses.ListUsersResponse{}
	res, err := c.httpGet("user/list", requestJSON(), requestExpects(200))
	if err != nil {
		return nil, err
	}
	if jsonErr := json.Unmarshal(res, &listusers); jsonErr != nil {
		return nil, &UnmarshalError{msg: multierror.Append(errDecoding, jsonErr).Error()}
	}
	for _, u := range listusers {
		nu := User{Login: u.Login, FirstName: u.FirstName, LastName: u.LastName, Email: u.Email}
		// from api v21 to v27, new columns were added
		// we need to handle that here
		if u.Created != nil {
			nu.Created = u.Created.Time
		}
		if u.Updated != nil {
			nu.Updated = u.Updated.Time
		}
		if u.LastJob != nil {
			nu.LastJob = u.LastJob.Time
		}
		users = append(users, nu)
	}
	return users, nil
}

// GetCurrentUserProfile returns information about the current user
// http://rundeck.org/docs/api/index.html#get-user-profile
func (c *Client) GetCurrentUserProfile() (*User, error) {
	if err := c.checkRequiredAPIVersion(responses.UserProfileResponse{}); err != nil {
		return nil, err
	}
	user := &User{}
	res, err := c.httpGet("user/info", requestJSON(), requestExpects(200))
	if err != nil {
		return nil, err
	}
	if jsonErr := json.Unmarshal(res, &user); jsonErr != nil {
		return nil, &UnmarshalError{msg: multierror.Append(errDecoding, jsonErr).Error()}
	}
	return user, nil
}

// GetUserProfile returns information about the named user - requires admin privileges
// http://rundeck.org/docs/api/index.html#get-another-user-profile
func (c *Client) GetUserProfile(login string) (*User, error) {
	if err := c.checkRequiredAPIVersion(responses.UserProfileResponse{}); err != nil {
		return nil, err
	}
	user := &User{}
	res, err := c.httpGet("user/info/"+login, requestJSON(), requestExpects(200))
	if err != nil {
		return nil, err
	}
	if jsonErr := json.Unmarshal(res, &user); jsonErr != nil {
		return nil, &UnmarshalError{msg: multierror.Append(errDecoding, jsonErr).Error()}
	}
	return user, nil
}

// ModifyUserProfile updates a user
// http://rundeck.org/docs/api/index.html#modify-user-profile
func (c *Client) ModifyUserProfile(u *User) (*User, error) {
	if err := c.checkRequiredAPIVersion(responses.UserProfileResponse{}); err != nil {
		return nil, err
	}
	currentUser, currentUserErr := c.GetCurrentUserProfile()
	if currentUserErr != nil {
		return nil, currentUserErr
	}
	if u.Login == "nil" {
		return nil, errors.New("must provide login and at least one field to update")
	}
	updatePath := "user/info"
	if currentUser.Login != u.Login {
		updatePath = "user/info/" + u.Login
	}
	newUser := requests.UserInfo{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
	postData, postDataErr := json.Marshal(newUser)
	if postDataErr != nil {
		return nil, postDataErr
	}
	res, resErr := c.httpPost(updatePath, withBody(bytes.NewReader(postData)), requestJSON(), requestExpects(200))
	if resErr != nil {
		return nil, resErr
	}
	resUser := &User{}
	if jsonErr := json.Unmarshal(res, &resUser); jsonErr != nil {
		return nil, &UnmarshalError{msg: multierror.Append(errDecoding, jsonErr).Error()}
	}
	return resUser, nil
}
