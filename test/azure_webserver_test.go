package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// subscriptionID represents the Azure subscription ID.
var subscriptionID = "9f7141c2-48ab-4f22-bc4f-e9ea835b1ff8"

// TestAzureLinuxVMCreation tests the creation of an Azure Linux VM.
func TestAzureLinuxVMCreation(t *testing.T) {
	// Set up Terraform options.
	terraformOptions := &terraform.Options{
		TerraformDir: "../", // Terraform configuration directory.
		Vars: map[string]interface{}{
			"labelPrefix": "dewa0117", // Custom label prefix.
		},
	}

	defer terraform.Destroy(t, terraformOptions) // Clean up after test.

	terraform.InitAndApply(t, terraformOptions) // Apply Terraform configuration.

	vmName := terraform.Output(t, terraformOptions, "vm_name")                        // Get VM name.
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name") // Get resource group.

	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID)) // Check VM existence.
}

// TestAzureLinuxVMUbuntuVersion verifies if the deployed Azure Linux VM is running the specified Ubuntu version.
func TestAzureLinuxVMUbuntuVersion(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"labelPrefix": "dewa0117",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	expectedVmImageVersion := terraform.Output(t, terraformOptions, "vm_image_version")

	getVirtualMachineImage := azure.GetVirtualMachineImage(t, vmName, resourceGroupName, subscriptionID)

	assert.Equal(t, expectedVmImageVersion, getVirtualMachineImage.Version)
}

// TestAzureNicExistsAndConnectedVm checks if a Network Interface Card (NIC) exists and is properly attached to the specified VM.
func TestAzureNicExistsAndConnectedVm(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"labelPrefix": "dewa0117",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")

	listNic := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)

	assert.Contains(t, listNic, nicName)
}
