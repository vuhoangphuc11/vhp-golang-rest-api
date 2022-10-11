package data

import "github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/dto"

var Users []dto.User

func init() {
	Users = []dto.User{
		{ID: 1, Email: "tai.nguyen@gmail.com", FullName: "Nguyen Huu Tai", Password: "1234", Age: 27, Gender: true, Phone: "0969999999", IsActive: true, Role: "User"},
		{ID: 2, Email: "hung.pham@gmail.com", FullName: "Nguyen Manh Hung", Password: "1234", Age: 22, Gender: true, Phone: "0969999999", IsActive: true, Role: "User"},
		{ID: 3, Email: "hong.kim@gmail.com", FullName: "Kim Xuan Hong", Password: "1234", Age: 26, Gender: false, Phone: "0969999999", IsActive: true, Role: "User"},
		{ID: 4, Email: "phuc.vu@gmail.com", FullName: "Vu Hoang Phuc", Password: "1234", Age: 21, Gender: true, Phone: "0969999999", IsActive: true, Role: "Admin"},
		{ID: 5, Email: "hien.bui@gmail.com", FullName: "Bui Phi Hien", Password: "1234", Age: 22, Gender: false, Phone: "0969999999", IsActive: false, Role: "User"},
		{ID: 6, Email: "ngan.doan@gmail.com", FullName: "Doan Huynh Long Ngan", Password: "1234", Age: 22, Gender: false, Phone: "0969999999", IsActive: false, Role: "User"},
		{ID: 7, Email: "uy.huu@gmail.com", FullName: "Nguyen Huu Uy", Password: "1234", Age: 27, Gender: true, Phone: "0969999999", IsActive: true, Role: "User"},
		{ID: 8, Email: "nhi.nguyen@gmail.com", FullName: "Nguyen Thi Yen Nhi", Password: "1234", Age: 26, Gender: false, Phone: "0969999999", IsActive: true, Role: "User"},
		{ID: 9, Email: "lam.nguyen@gmail.com", FullName: "Nguyen Ngoc Lam", Password: "1234", Age: 35, Gender: true, Phone: "0969999999", IsActive: false, Role: "Admin"},
		{ID: 10, Email: "tri.tran@gmail.com", FullName: "Tran Trong Tri", Password: "1234", Age: 27, Gender: true, Phone: "0969999999", IsActive: true, Role: "User"},
		{ID: 11, Email: "anh.nguyen@gmail.com", FullName: "Nguyen Hoang Anh", Password: "1234", Age: 23, Gender: true, Phone: "0969999999", IsActive: true, Role: "User"},
	}
}
