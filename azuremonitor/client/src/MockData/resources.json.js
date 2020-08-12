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
                        usage: "Storage Being Used: 24%",
                        consumption: 20000,
                        savings: "-",
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
                        usage: "58gb",
                        consumption: 30000,
                        savings: "-",
                        recommendations: [
                            "On average, we are logging 58 GB per day and are currently on the 200 GB per day plan. You could save $xx.xx/day if you lower your capacity.",
                            "On average, you are using 3203 notifications per month. If you switch to the pay as you go plan, you could save $xx.xxx"
                        ]
                    },
                    {
                        resourceName: "VMInsights(DefaultWorkspace-bb07e91d-a908-4fe4-a04e-40cf2d4b0603-EUS)",
                        type: "Solution",
                        usage: "usage description 4",
                        consumption: 40010,
                        savings: "-",
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
                        usage: "usage description 3",
                        consumption: 30000,
                        savings: "-",
                        recommendations: [
                            "Recommendation #3",
                            "Recommendation #3.1",
                            "Recommendation #3.2"
                        ]
                    },
                    {
                        resourceName: "elysium",
                        type: "Virtual Machine",
                        usage: "usage description 4",
                        consumption: 40010,
                        savings: "-",
                        recommendations: [
                            "On average, only 32% of CPU is being used. You could save $xx.xx if you lower XXXX to XXXX",
                            "Only 16% of storage is being used. You could save $xx.xxx if you lower storage to XXXX"
                        ]
                    },
                    {
                        resourceName: "elysium-ip",
                        type: "Public IP Address",
                        usage: "usage description 5",
                        consumption: 50000,
                        savings: "-",
                        recommendations: [
                            
                        ]
                    },
                    {
                        resourceName: "elysium-nsg",
                        type: "Network Security Group",
                        usage: "usage description 5",
                        consumption: 50000,
                        savings: "-",
                        recommendations: [
                           "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam",
                           "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam",
                           "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, "
                        ]
                    },
                    {
                        resourceName: "elysium653",
                        type: "Network Interface",
                        usage: "usage description 5",
                        consumption: 50000,
                        savings: "-",
                        recommendations: [
                            "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam",
                            "Recommendation #5.1"
                        ]
                    },
                    {
                        resourceName: "elysium_demo-vnet",
                        type: "Virtual network",
                        usage: "usage description 5",
                        consumption: 50000,
                        savings: "-",
                        recommendations: [
                            "Recommendation #5",
                            "Recommendation #5.1"
                        ]
                    },
                    {
                        resourceName: "elysium_disk1_557a7572648b494bb31b135726187544",
                        type: "Disk",
                        usage: "usage description 5",
                        consumption: 50000,
                        savings: "-",
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