import React from 'react'
import {AppContext} from '../contexts/AppContext'
import {Tooltip, IconButton, Menu, MenuItem} from '@material-ui/core'
import MoreVertIcon from '@material-ui/icons/MoreVert';


const MoreVert = (props) => {
    const { state, dispatch } = React.useContext(AppContext)
    const [anchorEl, setAnchorEl] = React.useState(null)

    const handleClick = (event) => {
        setAnchorEl(event.currentTarget)
    }

    const handleClose = () => {
        setAnchorEl(null)
    }

    const handleUpdateSavings = (amount) => {
        const {subscription, groupName, resourceName} = props
        const payload = {
            subscription: subscription,
            group: groupName,
            resource: resourceName,
            amount: amount
        }

        dispatch({type: 'SET_SAVINGS', payload: payload})
        handleClose()
    }

    return (
        <div>
            <Tooltip title='Options'>
                <IconButton onClick={handleClick}>
                    <MoreVertIcon/>
                </IconButton>
            </Tooltip>

            <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleClose}
            >
                <MenuItem onClick={() => handleUpdateSavings('-')}>Premium SSD - P10</MenuItem>
                <MenuItem onClick={() => handleUpdateSavings(835.32)}>Standard SSD - E2</MenuItem>
            </Menu>
        </div>
        
    )
}

export default MoreVert