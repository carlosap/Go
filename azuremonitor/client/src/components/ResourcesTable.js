import React, {useState} from 'react'

// Import Material UI Components and Icons
import {withStyles} from '@material-ui/core/styles'
import {Table, TableHead, TableBody, TableRow, TableCell} from '@material-ui/core'
import {Box, IconButton, Collapse, Paper, Typography} from '@material-ui/core'
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';
import SubdirectoryArrowRightIcon from '@material-ui/icons/SubdirectoryArrowRight';

// Import Components
import Recommendations from './Recommendations'
import Edit from "./Edit"
import Usages from "./Usages"
import MoreVert from "./MoreVert"

const BorderlessCell = withStyles({
  root: {
    borderBottom: 'none'
  }
})(TableCell)

// Component for row with collapse
const CollapseRow = (props) => {
  const {name, open, setOpen, styles, dispatch} = props

  const handleClick = () => {
    if(open) {
      setOpen(false)
      dispatch({type: 'REMOVE_FROM_TABLE_STATE', payload: name})
    } else {
      setOpen(true)
      dispatch({type: 'ADD_TO_TABLE_STATE', payload: name})
    }
  }

  return (
    <TableRow>
      <BorderlessCell style={styles}>
          {name}
          <IconButton size='small' onClick={handleClick}>
            {open ? <KeyboardArrowUpIcon/> : <KeyboardArrowDownIcon/> }
          </IconButton>
      </BorderlessCell>
    </TableRow>
  )
}

// Wrapper Component for anything within resource groups
const ResourceGroups = (props) => {
  const {group, subscription, getInitialState, dispatch} = props
  let initialState = getInitialState(group.groupName)
  const [open, setOpen] = useState(initialState)
  const styles = {
    padding: {
      paddingBottom: '0', 
      paddingTop:'0',
    },
    collapseRow: {
      paddingBottom: '8px', 
      paddingTop:'8px',
      paddingLeft:'36px'
    }
  }
  const headerTexts = ['Resource Name', 'Type', 'Product', 'Consumption', 'Savings', 'Actions']
  
  return (
    <React.Fragment>
      <CollapseRow 
        style={{paddingLeft:'36px'}} 
        name={group.groupName} 
        open={open} 
        setOpen={setOpen} 
        styles={styles.collapseRow}
        dispatch={dispatch}
      />
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0, paddingLeft:'36px' }}>
          <Collapse in={open}> 
            <Box>
              {group.resources.length > 0 ? 
                <Table>
                  <TableHead>
                      <TableRow>
                        {headerTexts.map((text, idx) => {
                          // Add Icon and make make flex display for first header column
                          if(idx === 0) { 
                            return (
                              <TableCell key={text} style={{ display:'flex', paddingBottom: '8px', paddingTop: 0}}>
                                <SubdirectoryArrowRightIcon fontSize="small"/>
                                <Typography style={{fontWeight:"bold"}} variant="subtitle2">
                                  {text}
                                </Typography>
                              </TableCell>
                            )
                          } else return (
                            <TableCell align='center' key={text} style={{paddingBottom: '8px', paddingTop: 0}}>
                              <Typography style={{fontWeight:"bold"}} variant="subtitle2">
                                {text}
                              </Typography>
                            </TableCell>
                          )
                        })}
                      </TableRow>
                  </TableHead>
                  <TableBody>
                      {group.resources.map((resource, idx) => {
                        if(idx === group.resources.length - 1){
                          return (
                              <TableRow key={idx}>
                                <BorderlessCell style={styles.padding}>{resource.resourceName}</BorderlessCell>
                                <BorderlessCell align='center' style={styles.padding}>{resource.type}</BorderlessCell>
                                <BorderlessCell align='center' style={styles.padding}>
                                  {resource.product ? resource.product : "-"}
                                </BorderlessCell>
                                <BorderlessCell align='center' style={styles.padding}>${resource.consumption}</BorderlessCell>
                                <BorderlessCell align='center' style={styles.padding}>{resource.savings}</BorderlessCell>
                                <BorderlessCell align='center' style={styles.padding}>
                                  <div style={{display:'flex', justifyContent:'center'}}>
                                    <Usages resource={resource}/>
                                    <Recommendations recommendations={resource.recommendations}/>
                                    <Edit groupName={group.groupName} subscription={subscription} resource={resource}/>
                                    <MoreVert subscription={subscription} groupName={group.groupName} resourceName={resource.resourceName}/>
                                  </div>
                                </BorderlessCell>
                              </TableRow>
                          )
                        } else return (
                          <TableRow key={idx}>
                              <TableCell style={styles.padding}>{resource.resourceName}</TableCell>
                              <TableCell align='center' style={styles.padding}>{resource.type}</TableCell>
                              <TableCell align='center' style={styles.padding}>
                                {resource.product ? resource.product : "-"}
                              </TableCell>
                              <TableCell align='center' style={styles.padding}>${resource.consumption}</TableCell>
                              <TableCell align='center' style={styles.padding}>{resource.savings}</TableCell>
                              <TableCell align='center' style={styles.padding}>
                                <div style={{display:'flex', justifyContent:'center'}}>
                                  <Usages resource={resource}/>
                                  <Recommendations recommendations={resource.recommendations}/>
                                  <Edit groupName={group.groupName} subscription={subscription} resource={resource}/>
                                  <MoreVert subscription={subscription} groupName={group.groupName} resourceName={resource.resourceName}/>
                                </div>															
                              </TableCell>
                          </TableRow>
                        )
                      })}
                  </TableBody>
                </Table> : 
                <Typography> No Resources Found</Typography>
              }
            </Box>
          </Collapse>
        </TableCell>
      </TableRow>
    </React.Fragment>
  )
}

// Wrapper Component for anything within subscriptions
const Subscriptions = (props) => {

  const {subscription, getInitialState, dispatch} = props
  let initialState = getInitialState(subscription.subscriptionName)
  const [open, setOpen] = useState(initialState)
  const styles = {paddingBottom: '8px', paddingTop: '8px'}
  return (
    <React.Fragment>
      <CollapseRow 
        name={subscription.subscriptionName} 
        open={open} setOpen={setOpen} 
        styles={styles}
        dispatch={dispatch}
      />
      <TableRow>
          <TableCell style={{ paddingBottom: 0, paddingTop: 0 }}>
            <Collapse in={open}> 
              <Box>
                <Table>
                  <TableHead>
                    <TableRow>
                        <TableCell style={{display:'flex', paddingBottom: '8px', paddingTop:'0'}}>
                          <SubdirectoryArrowRightIcon fontSize="small"/>
                          <Typography style={{fontWeight:"bold"}} variant="subtitle2">
                            Resource Groups
                          </Typography>
                        </TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {subscription.resourceGroups.map((group,i) => (
                      <ResourceGroups 
                        subscription={subscription.subscriptionName} 
                        key={i} 
                        group={group}
                        getInitialState={getInitialState}
                        dispatch={dispatch}
                      />
                    ))}
                  </TableBody>
                </Table>
              </Box>
            </Collapse>
          </TableCell>
      </TableRow>
    </React.Fragment>
  )
}

const ResourcesTable = (props) => {
  const {data, tableState, dispatch} = props

  const getInitialState = (id) => { 
    return tableState.includes(id)
  }

  return (
    <Paper elevation={3}>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell> 
                <Typography style={{fontWeight:"bold"}} variant="h6">
                  Subscriptions
                </Typography> 
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {data.map((subscription, idx) => (
            <Subscriptions 
              key={subscription.subscriptionName} 
              subscription={subscription}
              tableState={tableState}
              getInitialState={getInitialState}
              dispatch={dispatch}
            />
          ))}
        </TableBody>
      </Table>
    </Paper>
   )
}

export default ResourcesTable