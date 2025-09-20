package checkIn

import (
	"fmt"
	"testing"
)

func TestIsOver18WhenConStarts(t *testing.T) {
	// ❶ Arrange
	expected := true

	//Con starts 2025-10-10
	bornYear := 2025 - 18
	born := fmt.Sprintf("%d-10-10", bornYear)

	// ❷ Act
	result := isOver18(born)

	// ❸ Assert
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestIsNotOver18WhenConStarts(t *testing.T) {
	// ❶ Arrange
	expected := false

	//Con starts 2025-10-10
	bornYear := 2025 - 18
	born := fmt.Sprintf("%d-10-11", bornYear)

	// ❷ Act
	result := isOver18(born)

	// ❸ Assert
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
