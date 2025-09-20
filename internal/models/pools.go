package models

import (
	"sync"
)

var (
	classifierPool = sync.Pool{
		New: func() interface{} {
			return &Classifier{}
		},
	}
)

func getClassifier() *Classifier {
	return classifierPool.Get().(*Classifier)
}

func putClassifier(c *Classifier) {
	c.ID = 0
	c.Name = ""
	c.CreatedAt = c.CreatedAt.UTC()
	classifierPool.Put(c)
}