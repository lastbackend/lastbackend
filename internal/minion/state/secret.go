//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package state

import (
	"context"
	"github.com/lastbackend/lastbackend/tools/logger"
	"sync"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

const logSecretPrefix = "state:secret:>"

type SecretsState struct {
	lock    sync.RWMutex
	secrets map[string]models.Secret
}

func (s *SecretsState) GetSecrets() map[string]models.Secret {
	log := logger.WithContext(context.Background())
	log.Debugf("%s get pods", logSecretPrefix)
	return s.secrets
}

func (s *SecretsState) SetSecrets(secrets map[string]*models.Secret) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s set secrets: %d", logSecretPrefix, len(secrets))
	for h, secret := range secrets {
		s.secrets[h] = *secret
	}
}

func (s *SecretsState) GetSecret(name string) *models.Secret {
	log := logger.WithContext(context.Background())
	log.Debugf("%s get secret: %s", logSecretPrefix, name)
	s.lock.Lock()
	defer s.lock.Unlock()
	pod, ok := s.secrets[name]
	if !ok {
		return nil
	}
	return &pod
}

func (s *SecretsState) AddSecret(name string, secret *models.Secret) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s add secret: %s", logSecretPrefix, name)
	s.SetSecret(name, secret)
}

func (s *SecretsState) SetSecret(name string, secret *models.Secret) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s set secret: %s", logSecretPrefix, name)
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.secrets[secret.SelfLink().String()]; ok {
		delete(s.secrets, secret.SelfLink().String())
	}

	s.secrets[secret.GetHash()] = *secret
}

func (s *SecretsState) DelSecret(name string) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s del secret: %s", logSecretPrefix, name)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.secrets[name]; ok {
		delete(s.secrets, name)
	}
}
