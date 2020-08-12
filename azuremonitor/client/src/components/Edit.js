import React from 'react'
import {withRouter} from 'react-router-dom'
import {IconButton, Tooltip} from '@material-ui/core'
import EditIcon from '@material-ui/icons/Edit';

const Edit = (props) => {
    const {resource, subscription} = props
    const handleClick = () => {
        props.history.push({
            pathname: `/analysis/${resource.type}/${resource.resourceName}`,
            state: { subscription: subscription}
        })
    }

    return (
        <Tooltip title="Analysis Form">
            <IconButton onClick={handleClick}>
                <EditIcon color="primary"/>
            </IconButton>
        </Tooltip>
    )
}

export default withRouter(Edit)