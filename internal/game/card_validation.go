package game

import "fmt"

func ValidateCardRegistry() error {
	if len(CardRegistry) == 0 {
		return fmt.Errorf("card registry is empty")
	}

	for id, definition := range CardRegistry {
		if definition.ID == "" {
			return fmt.Errorf("card %q has empty ID", id)
		}

		if definition.ID != id {
			return fmt.Errorf("card registry key %q does not match definition ID %q", id, definition.ID)
		}

		if definition.Name == "" {
			return fmt.Errorf("card %q has empty name", id)
		}

		if definition.Type == "" {
			return fmt.Errorf("card %q has empty type", id)
		}
	}

	return nil
}

func ValidateDeckCardsExist(deck []CardInstance) error {
	for _, card := range deck {
		if _, ok := CardRegistry[card.DefinitionID]; !ok {
			return fmt.Errorf("unknown card definition %q", card.DefinitionID)
		}
	}

	return nil
}
