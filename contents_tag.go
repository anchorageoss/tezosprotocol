package tezosprotocol

// ContentsTag captures the possible tag values for operation contents
type ContentsTag byte

const (
	// ContentsTagRevelation is the tag for revelations
	ContentsTagRevelation ContentsTag = 107
	// ContentsTagTransaction is the tag for transactions
	ContentsTagTransaction ContentsTag = 108
	// ContentsTagOrigination is the tag for originations
	ContentsTagOrigination ContentsTag = 109
	// ContentsTagDelegation is the tag for delegations
	ContentsTagDelegation ContentsTag = 110
	// ContentsTagEndorsement is the tag for endorsements
	ContentsTagEndorsement ContentsTag = 0
)
