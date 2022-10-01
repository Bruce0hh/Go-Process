package abstract

import (
	"fmt"
	"go-designpattern/factory2/abstract/adidas"
	"go-designpattern/factory2/abstract/nike"
	"go-designpattern/factory2/abstract/sportfactory"
)

var sportsFactoryMap = map[string]sportfactory.ISportsFactory{
	"adidas": &adidas.Adidas{},
	"nike":   &nike.Nike{},
}

func GetSportsFactory(brand string) (sportfactory.ISportsFactory, error) {
	if f, ok := sportsFactoryMap[brand]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("wrong brand type passed")
}
