const data = 
[
    {
        subscriptionName: "SpartanAppSolutions",
        resourceGroups: [
            {
                groupName: "cloud-shell-storage-eastus",
                resources: [
                    {
                        resourceName: "cs210032000a0f50565",
                        type: "Storage Account",
                        consumption: "Consumption description 1",
                        usage: 20000,
                        savings: 5000,
                        recommendations: [
                            "Recommendation #1",
                            "Recommendation #1.1",
                            "Recommendation #1.2",
                        ]
                    }
                ]
            },
            {
                groupName: "DefaultResourceGroup-EUS",
                resources: [
                    {
                        resourceName: "DefaultWorkspace-bb07e91d-a908-4fe4-a04e-40cf2d4b0603-EUS",
                        type: "Log Analytics workspace",
                        consumption: "Consumption description 3",
                        usage: 30000,
                        savings: 1000,
                        recommendations: [
                            "Recommendation #3",
                            "Recommendation #3.1",
                            "Recommendation #3.2"
                        ]
                    },
                    {
                        resourceName: "VMInsights(DefaultWorkspace-bb07e91d-a908-4fe4-a04e-40cf2d4b0603-EUS)",
                        type: "Solution",
                        consumption: "Consumption description 4",
                        usage: 40010,
                        savings: 3000,
                        recommendations: [
                            "Recommendation #4"
                        ]
                    },
                ]
            },
            {
                groupName: "elysium_demo",
                resources: [
                    {
                        resourceName: "edaemonnews",
                        type: "Cognitive Service",
                        consumption: "Consumption description 3",
                        usage: 30000,
                        savings: 1000,
                        recommendations: [
                            "Recommendation #3",
                            "Recommendation #3.1",
                            "Recommendation #3.2"
                        ]
                    },
                    {
                        resourceName: "elysium",
                        type: "Virtual Machine",
                        consumption: "Consumption description 4",
                        usage: 40010,
                        savings: 3000,
                        recommendations: [
                            "Recommendation #4"
                        ]
                    },
                    {
                        resourceName: "elysium-ip",
                        type: "Public IP Address",
                        consumption: "Consumption description 5",
                        usage: 50000,
                        savings: 5000,
                        recommendations: [
                            "Recommendation #5",
                            "Recommendation #5.1"
                        ]
                    },
                    {
                        resourceName: "elysium-nsg",
                        type: "Network Security Group",
                        consumption: "Consumption description 5",
                        usage: 50000,
                        savings: 5000,
                        recommendations: [
                            "Recommendation #5",
                            "Recommendation #5.1"
                        ]
                    },
                    {
                        resourceName: "elysium-nsg",
                        type: "Network Security Group",
                        consumption: "Consumption description 5",
                        usage: 50000,
                        savings: 5000,
                        recommendations: [
                            "Recommendation #5",
                            "Recommendation #5.1"
                        ]
                    },
                    {
                        resourceName: "elysium653",
                        type: "Network Interface",
                        consumption: "Consumption description 5",
                        usage: 50000,
                        savings: 5000,
                        recommendations: [
                            "Recommendation #5",
                            "Recommendation #5.1"
                        ]
                    },
                    {
                        resourceName: "elysium_demo-vnet",
                        type: "Virtual network",
                        consumption: "Consumption description 5",
                        usage: 50000,
                        savings: 5000,
                        recommendations: [
                            "Recommendation #5",
                            "Recommendation #5.1"
                        ]
                    },
                    {
                        resourceName: "elysium_disk1_557a7572648b494bb31b135726187544",
                        type: "Disk",
                        consumption: "Consumption description 5",
                        usage: 50000,
                        savings: 5000,
                        recommendations: [
                            "Recommendation #5",
                            "Recommendation #5.1"
                        ]
                    },
                ]
                
            },
            {
                groupName: "NetworkWatcherRG",
                resources: [
                   
                ]
            }
        ]
    }
]


export default data