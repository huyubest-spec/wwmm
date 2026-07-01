package controller

import (
	"encoding/json"
	"wwmm/dao"
)

func prettyJSON(s string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return s
	}
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b)
}

func hasVotedDB(uid, pid int) (bool, error) {
	return dao.HasVoted(uid, pid)
}
