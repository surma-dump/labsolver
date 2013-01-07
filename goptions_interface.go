package main

import (
	"encoding/json"
	"fmt"
)

func (c *Crop) MarshalGoption(s string) error {
	values := []int{}
	err := json.Unmarshal([]byte(s), &values)
	if err != nil {
		return err
	}
	if len(values) != 4 {
		return fmt.Errorf("Expected array of length 4, got %d", len(values))
	}
	copy([]int(c[:]), values)
	return nil
}

func (c *Crop) String() string {
	return fmt.Sprintf("[%d, %d, %d, %d]", c[0], c[1], c[2], c[3])
}

func (p *Vector2) MarshalGoption(s string) error {
	values := []int{}
	err := json.Unmarshal([]byte(s), &values)
	if err != nil {
		return err
	}
	if len(values) != 2 {
		return fmt.Errorf("Expected array of length 2, got %d", len(values))
	}
	p.Point.X, p.Point.Y = values[0], values[1]
	return nil
}

func (p *Vector2) String() string {
	return fmt.Sprintf("[%d, %d]", p.X, p.Y)
}
