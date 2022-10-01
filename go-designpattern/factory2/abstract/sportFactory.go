package abstract

import (
	"fmt"
	"go-designpattern/factory2/abstract/adidas"
	"go-designpattern/factory2/abstract/nike"
	"go-designpattern/factory2/abstract/sportfactory"
)

func GetSportsFactory(brand string) (sportfactory.ISportsFactory, error) {
	if brand == "adidas" {
		return &adidas.Adidas{}, nil
	}

	if brand == "nike" {
		return &nike.Nike{}, nil
	}

	return nil, fmt.Errorf("wrong brand type passed")
}
