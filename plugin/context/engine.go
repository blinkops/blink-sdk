package context

import (
	"errors"
	"fmt"
	"strings"
)

func getMap(parameter string, context map[string]interface{}) map[string]interface{} {
	// Usage: parameter will be path (json.dot.walking) that the start of the parameter
	// 		  strings.Split(parameter, ".")[0] -> Should be `inputs` or `parameters`
	if !strings.Contains(parameter, ".") {
		return nil
	}

	location := strings.Split(parameter, ".")[0]

	internalContext, ok := context[location]

	if !ok {
		return nil
	}

	return internalContext.(map[string] interface{})
}

func getIterator(parameter string, create bool, context map[string]interface{}) (map[string]interface{}, error) {
	if !strings.Contains(parameter, ".") {
		return nil, fmt.Errorf("provided parameter has to be path, seperated by '.', %v", parameter)
	}

	location := strings.Split(parameter, ".")
	depth := len(location) - 1
	if depth == 0 || location[1] == "" {
		return nil, fmt.Errorf("provided parameter is not allowed, must be at-least 1 depth, %v", parameter)
	}

	iterateMap := getMap(parameter, context)
	if iterateMap == nil {
		return nil, fmt.Errorf("failed to get iterator with parameter: %v", parameter)
	}

	iterator := interface{}(iterateMap)
	for i := 1; i < depth; i++ {
		iteratorMap, ok := iterator.(map[string]interface{})
		if !ok {
			return nil, errors.New("failed to convert iterator to map[string] interface")
		}
		currentHead := location[i]
		switch iteratorMap[currentHead].(type) {
		case map[string]interface{}:
			break
		default:
			if create {
				iteratorMap[currentHead] = make(map[string]interface{})
			} else {
				return nil, errors.New("given path doesn't exists")
			}
		}
		iterator = iteratorMap[currentHead]
	}
	return iterator.(map[string]interface{}), nil
}

func Get(parameter string, context map[string]interface{}) (interface{}, error) {
	iterator, err := getIterator(parameter, false, context)
	if err != nil || iterator == nil {
		return nil, err
	}
	location := strings.Split(parameter, ".")
	return iterator[location[len(location)-1]], nil
}

func Set(parameter string, value interface{}, context map[string]interface{}) error {
	iterator, err := getIterator(parameter, true, context)
	if err != nil || iterator == nil {
		return err
	}
	location := strings.Split(parameter, ".")
	iterator[location[len(location)-1]] = value
	return nil
}

func Delete(parameter string, context map[string]interface{}) error {
	iterator, err := getIterator(parameter, true, context)
	if err != nil || iterator == nil {
		return err
	}
	location := strings.Split(parameter, ".")
	delete(iterator, location[len(location)-1])
	return nil
}

