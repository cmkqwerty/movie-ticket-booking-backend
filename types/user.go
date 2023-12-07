package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

const (
	bcryptCost        = 12
	minFirstNameLen   = 2
	minLastNameLen    = 2
	minPasswordLength = 8
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (p CreateUserParams) Validate() map[string]string {
	errs := map[string]string{}

	if len(p.FirstName) < minFirstNameLen {
		errs["firstName"] = fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen)
	}

	if len(p.LastName) < minLastNameLen {
		errs["lastName"] = fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen)
	}

	if len(p.Password) < minPasswordLength {
		errs["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLength)
	}

	if !isEmailValid(p.Email) {
		errs["email"] = fmt.Sprintf("invalid email address")
	}

	return errs
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`)

	return emailRegex.MatchString(email)
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (p UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}

	if len(p.FirstName) > 0 {
		m["firstName"] = p.FirstName
	}

	if len(p.LastName) > 0 {
		m["lastName"] = p.LastName
	}

	return m
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encPw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encPw),
	}, nil
}
