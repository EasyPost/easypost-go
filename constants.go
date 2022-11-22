package easypost

const FedExAccountType = "FedexAccount"
const UpsAccountType = "UpsAccount"

func getCustomWorkflowCarrierAccountTypes() []string {
	return []string{FedExAccountType, UpsAccountType}
}
