package room

func (r *room) autoAbandonSuit() {
	for i, user := range r.users {
		if user < 10000 && 0 == r.abandonSuits[i] {
			var dot int
			var bamboo int
			var character int
			for _, tile := range r.userTiles.HoldTiles[i] {
				switch tile.Suit {
				case suitDot:
					dot++
				case suitBamboo:
					bamboo++
				case suitCharacter:
					character++
				}
			}
			if dot < bamboo && dot < character {
				r.abandonSuits[i] = suitDot
			} else if bamboo < character {
				r.abandonSuits[i] = suitBamboo
			} else {
				r.abandonSuits[i] = suitCharacter
			}
		}
	}
}
