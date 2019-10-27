package sync

import (
	"fmt"

	"github.com/coveooss/credentials-sync/credentials"
	"github.com/coveooss/credentials-sync/targets"
	log "github.com/sirupsen/logrus"
)

// DeleteListOfCredentials deletes the configured list of credentials from the given target
func (config *Configuration) DeleteListOfCredentials(target targets.Target) error {
	for _, id := range config.CredentialsToDelete {
		if targets.HasCredential(target, id) {
			log.Infof("[%s] Deleting %s", target.GetName(), id)
			if err := target.DeleteCredentials(id); err != nil {
				err = fmt.Errorf("Failed to delete credentials with ID %s from %s: %v", id, target.GetName(), err)
				if config.StopOnError {
					return err
				}
				log.Error(err)
			}
		}
	}

	return nil
}

// UpdateListOfCredentials syncs the given list of credentials to the given target
func (config *Configuration) UpdateListOfCredentials(target targets.Target, listOfCredentials []credentials.Credentials) error {
	isSynced := func(id string) bool {
		for _, credentials := range listOfCredentials {
			if credentials.GetID() == id {
				return true
			}
		}
		return false
	}

	for _, credentials := range listOfCredentials {
		log.Infof("[%s] Syncing %s", target.GetName(), credentials.GetID())
		if err := target.UpdateCredentials(credentials); err != nil {
			err = fmt.Errorf("Failed to send credentials with ID %s to %s: %v", credentials.GetID(), target.GetName(), err)
			if config.StopOnError {
				return err
			}
			log.Error(err)
		}
	}

	if target.ShouldDeleteUnsynced() {
		log.Debugf("Deleting unsynced credentials from %v", target.GetName())
	}

	for _, existingID := range target.GetExistingCredentials() {
		if !isSynced(existingID) {
			if target.ShouldDeleteUnsynced() {
				log.Infof("[%s] Deleting %s", target.GetName(), existingID)
				if err := target.DeleteCredentials(existingID); err != nil {
					err = fmt.Errorf("Failed to delete credentials with ID %s from %s: %v", existingID, target.GetName(), err)
					if config.StopOnError {
						return err
					}
					log.Error(err)
				}
			} else {
				log.Infof("[%s] %s is unsynced. Not modifying it", target.GetName(), existingID)
			}
		}
	}

	return nil
}
