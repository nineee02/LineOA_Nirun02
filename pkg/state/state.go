package state

// import "sync"

// var userState sync.Map

// // SetUserState - ตั้งสถานะผู้ใช้
// func SetUserState(userID, state string) {
// 	userState.Store(userID, state)
// }

// // GetUserState - ดึงสถานะผู้ใช้
// func GetUserState(userID string) (string, bool) {
// 	state, exists := userState.Load(userID)
// 	if exists {
// 		return state.(string), exists
// 	}
// 	return "", exists
// }

// // ClearUserState - ลบสถานะผู้ใช้
// func ClearUserState(userID string) {
// 	userState.Delete(userID)
// }
