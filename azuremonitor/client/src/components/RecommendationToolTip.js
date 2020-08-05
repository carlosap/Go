import React from 'react'
import {Tooltip, IconButton} from '@material-ui/core'
import InfoIcon from '@material-ui/icons/Info';

const RecommendationToolTip = (props) => {    
    const createRecommendations = () => {
        return "test"
    }
    return (
        <Tooltip title={createRecommendations(props.recommendations)}>
            <InfoIcon color='primary'/>
        </Tooltip>
    )
}

export default RecommendationToolTip