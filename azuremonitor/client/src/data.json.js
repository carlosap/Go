export const data = {
    subscription: {
        name: "Subscription 1",
        resourceGroups: [
            {
                groupName: "Resource Group 1",
                resources: [
                    {
                        resourceName: "Group 1 Resource 1",
                        type: "VM",
                        consumptions: "Consumption description 1",
                        usage: 20000,
                        savings: 5000,
                        recommendations: "Reccommendation #1"
                    },
                    {
                        resourceName: "Group 1 Resource 2",
                        type: "Other",
                        consumptions: "Consumption description 2",
                        usage: 50000,
                        savings: 3000,
                        recommendations: "Reccommendation #2"
                    }
                ]
            },
            {
                groupName: "Resource Group 2",
                resources: [
                    {
                        resourceName: "Group 2 Resource 1",
                        type: "VM",
                        consumptions: "Consumption description 3",
                        usage: 30000,
                        savings: 1000,
                        recommendations: "Reccommendation #3"
                    },
                    {
                        resourceName: "Group 2 Resource 2",
                        type: "Other",
                        consumptions: "Consumption description 4",
                        usage: 40010,
                        savings: 3000,
                        recommendations: "Reccommendation #4"
                    },
                    {
                        resourceName: "Group 2 Resource 3",
                        type: "Other",
                        consumptions: "Consumption description 5",
                        usage: 50000,
                        savings: 5000,
                        recommendations: "Reccommendation #5"
                    },
                ]
            }
        ]
    }
}