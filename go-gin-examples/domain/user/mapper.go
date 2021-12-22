package user

func ToDTO(user *User) *UserDTO {
	if user == nil {
		return nil
	}

	return &UserDTO{
		Id:               user.Id.Hex(),
		Username:         user.Username,
		Roles:            user.Roles,
		Email:            user.Email,
		Activated:        user.Activated,
		ActivatedDate:    user.ActivatedDate.Unix(),
		AvatarUri:        user.AvatarUri,
		Describe:         user.Describe,
		CreatedDate:      user.CreatedDate,
		LastModifiedDate: user.LastModifiedDate,
	}
}

func ToDTOS(users []*User) []*UserDTO {
	if users == nil || len(users) == 0 {
		return make([]*UserDTO, 0)
	}

	dtos := make([]*UserDTO, len(users))
	for i, user := range users {
		dtos[i] = ToDTO(user)
	}
	return dtos
}
