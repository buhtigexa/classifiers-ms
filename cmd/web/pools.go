package main

import (
	"sync"

	"classifier.buhtigexa.net/internal/models"
)

var (
	classifierPool = sync.Pool{
		New: func() interface{} {
			return &models.Classifier{}
		},
	}

	envelopePool = sync.Pool{
		New: func() interface{} {
			return make(envelope, 4)
		},
	}

	// Pool para slices de clasificadores con una capacidad inicial típica
	classifierSlicePool = sync.Pool{
		New: func() interface{} {
			slice := make([]*models.Classifier, 0, 20) // capacidad típica para una página
			return &slice
		},
	}
)

func getClassifier() *models.Classifier {
	return classifierPool.Get().(*models.Classifier)
}

func putClassifier(c *models.Classifier) {
	c.ID = 0
	c.Name = ""
	c.CreatedAt = c.CreatedAt.UTC()
	classifierPool.Put(c)
}

func getEnvelope() envelope {
	env := envelopePool.Get().(envelope)
	for k := range env {
		delete(env, k)
	}
	return env
}

func putEnvelope(env envelope) {
	for k := range env {
		delete(env, k)
	}
	envelopePool.Put(env)
}

func getClassifierSlice() *[]*models.Classifier {
	return classifierSlicePool.Get().(*[]*models.Classifier)
}

func putClassifierSlice(s *[]*models.Classifier) {
	*s = (*s)[:0] // clear slice but keep capacity
	classifierSlicePool.Put(s)
}