package blockchain

import (
	"wwmm/utils"
)

func MerkleRoot(leafHashes []string) string {
	if len(leafHashes) == 0 {
		return utils.Sha256HexString("EMPTY_MERKLE_ROOT")
	}
	if len(leafHashes) == 1 {
		return leafHashes[0]
	}
	level := leafHashes
	for len(level) > 1 {
		var next []string
		for i := 0; i < len(level); i += 2 {
			if i+1 < len(level) {
				next = append(next, utils.Sha256HexString(level[i]+level[i+1]))
			} else {
				next = append(next, utils.Sha256HexString(level[i]+level[i]))
			}
		}
		level = next
	}
	return level[0]
}
