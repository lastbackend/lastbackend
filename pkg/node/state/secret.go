//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

func (s *SecretsState) GetSecrets() map[string]types.Secret {
	log.V(logLevel).Debug("Cache: SecretCache: get pods")
	return s.secrets
}

func (s *SecretsState) SetSecrets(secrets []*types.Secret) {
	log.V(logLevel).Debugf("Cache: SecretCache: set secrets: %#v", secrets)
	for _, secret := range secrets {
		s.secrets[secret.GetHash()] = *secret
	}
}

func (s *SecretsState) GetSecret(hash string) *types.Secret {
	log.V(logLevel).Debugf("Cache: SecretCache: get secret: %s", hash)
	s.lock.Lock()
	defer s.lock.Unlock()
	pod, ok := s.secrets[hash]
	if !ok {
		return nil
	}
	return &pod
}

func (s *SecretsState) AddSecret(secret *types.Secret) {
	log.V(logLevel).Debugf("Cache: SecretCache: add secret: %#v", secret)
	s.SetSecret(secret)
}

func (s *SecretsState) SetSecret(secret *types.Secret) {
	log.V(logLevel).Debugf("Cache: SecretCache: set secret: %#v", secret)
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.secrets[secret.GetHash()]; ok {
		delete(s.secrets, secret.GetHash())
	}

	s.secrets[secret.GetHash()] = *secret
}

func (s *SecretsState) DelSecret(secret *types.Secret) {
	log.V(logLevel).Debugf("Cache: SecretCache: del secret: %#v", secret)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.secrets[secret.GetHash()]; ok {
		delete(s.secrets, secret.GetHash())
	}
}
