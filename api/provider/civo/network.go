package civo

import (
	"fmt"

	"github.com/kubesimplify/ksctl/api/resources"
	"github.com/kubesimplify/ksctl/api/utils"
)

// NewNetwork implements resources.CloudInfrastructure.
func (obj *CivoProvider) NewNetwork(state resources.StateManagementInfrastructure) error {

	// check if the networkID already exist
	if len(civoCloudState.NetworkIDs.NetworkID) != 0 {
		fmt.Println("[skip] network creation found", civoCloudState.NetworkIDs.NetworkID)
		return nil // its a part of play back
	}

	res, err := civoClient.NewNetwork(obj.Metadata.ResName)
	if err != nil {
		return err
	}
	civoCloudState.NetworkIDs.NetworkID = res.ID
	fmt.Printf("[civo] Created %s network\n", obj.Metadata.ResName)

	// NOTE: as network creation marks first resource we should create the directoy
	// when its success

	if err := state.Path(generatePath(utils.CLUSTER_PATH, clusterType, clusterDirName)).
		Permission(FILE_PERM_CLUSTER_DIR).CreateDir(); err != nil {
		return err
	}

	path := generatePath(utils.CLUSTER_PATH, clusterType, clusterDirName, STATE_FILE_NAME)

	return saveStateHelper(state, path)
}

// DelNetwork implements resources.CloudInfrastructure.
func (obj *CivoProvider) DelNetwork(state resources.StateManagementInfrastructure) error {

	if len(civoCloudState.NetworkIDs.NetworkID) == 0 {
		fmt.Println("[skip] network already deleted")
	} else {
		_, err := civoClient.DeleteNetwork(civoCloudState.NetworkIDs.NetworkID)
		if err != nil {
			return err
		}
		fmt.Printf("[civo] Deleted %s network\n", civoCloudState.NetworkIDs.NetworkID)
		civoCloudState.NetworkIDs.NetworkID = ""
		if err := saveStateHelper(state, generatePath(utils.CLUSTER_PATH, clusterType, clusterDirName, STATE_FILE_NAME)); err != nil {
			return err
		}
	}
	path := generatePath(utils.CLUSTER_PATH, clusterType, clusterDirName)
	return state.Path(path).DeleteDir()
}
