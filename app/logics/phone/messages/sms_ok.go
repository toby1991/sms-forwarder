package messages


func ParseOk(msg []byte) bool {
	msgStr :=  string(msg[:])
	if msgStr == "OK"{
		return true
	}
	return false
}
