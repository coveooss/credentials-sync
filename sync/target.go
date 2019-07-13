package sync

import (
	"github.com/coveooss/credentials-sync/credentials"
	"github.com/coveooss/credentials-sync/targets"
	log "github.com/sirupsen/logrus"
)

// DeleteListOfCredentials deletes the configured list of credentials from the given target
func (config *Configuration) DeleteListOfCredentials(target targets.Target) {
	for _, id := range config.CredentialsToDelete {
		if targets.HasCredential(target, id) {
			log.Infof("[%s] Deleting %s", target.GetName(), id)
			if err := target.DeleteCredentials(id); err != nil {
				config.logError("Failed to delete credential %s from %s: %v", id, target.GetName(), err)
			}
		}
	}
}

// UpdateListOfCredentials syncs the given list of credentials to the given target
func (config *Configuration) UpdateListOfCredentials(target targets.Target, listOfCredentials []credentials.Credentials) {
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
			config.logError("Failed to send credentials with ID %s to %s: %v", credentials.GetID(), target.GetName(), err)
		}
	}

	if target.ShouldDeleteUnsynced() {
		log.Debugf("Deleting unsynced credentials from %v", target.GetName())
		for _, existingID := range target.GetExistingCredentials() {
			if !isSynced(existingID) {
				log.Infof("[%s] Deleting %s", target.GetName(), existingID)
				if err := target.DeleteCredentials(existingID); err != nil {
					config.logError("Failed to delete credentials with ID %s to %s: %v", existingID, target.GetName(), err)
				}
			}
		}
	}
}
