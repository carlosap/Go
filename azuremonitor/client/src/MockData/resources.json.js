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
                        consumption: "Storage Being Used: 24%",
                        usage: 20000,
                        savings: 4000,
                        recommendations: [
                            "Only 24% of your storage is being used. You could save $xxx.xx/month if you switch to a smaller max storage size.",
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
                        consumption: "58gb",
                        usage: 30000,
                        savings: 1000,
                        recommendations: [
                            "On average, we are logging 58 GB per day and are currently on the 200 GB per day plan. You could save $xx.xx/day if you lower your capacity.",
                            "On average, you are using 3203 notifications per month. If you switch to the pay as you go plan, you could save $xx.xxx"
                        ]
                    },
                    {
                        resourceName: "VMInsights(DefaultWorkspace-bb07e91d-a908-4fe4-a04e-40cf2d4b0603-EUS)",
                        type: "Solution",
                        consumption: "Consumption description 4",
                        usage: 40010,
                        savings: 3000,
                        recommendations: [
                            
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
                            "On average, only 32% of CPU is being used. You could save $xx.xx if you lower XXXX to XXXX",
                            "Only 16% of storage is being used. You could save $xx.xxx if you lower storage to XXXX"
                        ]
                    },
                    {
                        resourceName: "elysium-ip",
                        type: "Public IP Address",
                        consumption: "Consumption description 5",
                        usage: 50000,
                        savings: 5000,
                        recommendations: [
                            
                        ]
                    },
                    {
                        resourceName: "elysium-nsg",
                        type: "Network Security Group",
                        consumption: "Consumption description 5",
                        usage: 50000,
                        savings: 5000,
                        recommendations: [
                           "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam",
                           "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam",
                           "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, "
                        ]
                    },
                    {
                        resourceName: "elysium653",
                        type: "Network Interface",
                        consumption: "Consumption description 5",
                        usage: 50000,
                        savings: 5000,
                        recommendations: [
                            "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam",
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