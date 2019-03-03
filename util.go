package main

func string2bool(s string) bool {
	if s == "true" || s == "TRUE" {
		return true
	}
	return false

}
