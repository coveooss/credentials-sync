package sync

import (
	"fmt"

	"github.com/coveooss/credentials-sync/credentials"
	"github.com/coveooss/credentials-sync/logger"
	"github.com/coveooss/credentials-sync/targets"
)

// DeleteListOfCredentials deletes the configured list of credentials from the given target
func (config *Configuration) DeleteListOfCredentials(target targets.Target) error {
	for _, id := range config.CredentialsToDelete {
		if targets.HasCredential(target, id) {
			logger.Log.Infof("[%s] Deleting %s", target.GetName(), id)
			if err := target.DeleteCredentials(id); err != nil {
				err = fmt.Errorf("Failed to delete credentials with ID %s from %s: %v", id, target.GetName(), err)
				if config.StopOnError {
					return err
				}
				logger.Log.Error(err)
			}
		}
	}

	return nil
}

// UpdateListOfCredentials syncs the given list of credentials to the given target
func (config *Configuration) UpdateListOfCredentials(target targets.Target, listOfCredentials []credentials.Credentials) error {
	isSynced := func(id string) bool {
		for _, credentials := range listOfCredentials {
			if credentials.GetTargetID() == id {
				return true
			}
		}
		return false
	}

	for _, credentials := range listOfCredentials {
		logger.Log.Infof("[%s] Syncing %s", target.GetName(), credentials.GetTargetID())
		if err := target.UpdateCredentials(credentials); err != nil {
			err = fmt.Errorf("Failed to send credentials with ID %s to %s: %v", credentials.GetTargetID(), target.GetName(), err)
			if config.StopOnError {
				return err
			}
			logger.Log.Error(err)
		}
	}

	if target.ShouldDeleteUnsynced() {
		logger.Log.Debugf("Deleting unsynced credentials from %v", target.GetName())
	}

	for _, existingID := range target.GetExistingCredentials() {
		if !isSynced(existingID) {
			if target.ShouldDeleteUnsynced() {
				logger.Log.Infof("[%s] Deleting %s", target.GetName(), existingID)
				if err := target.DeleteCredentials(existingID); err != nil {
					err = fmt.Errorf("Failed to delete credentials with ID %s from %s: %v", existingID, target.GetName(), err)
					if config.StopOnError {
						return err
					}
					logger.Log.Error(err)
				}
			} else {
				logger.Log.Infof("[%s] %s is unsynced. Not modifying it", target.GetName(), existingID)
			}
		}
	}

	return nil
}
