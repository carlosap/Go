import React, {useState} from 'react'
import {makeStyles} from '@material-ui/core/styles'
import {Table, TableHead, TableBody, TableRow, TableCell} from '@material-ui/core'
import {Box, IconButton, Collapse, Paper, Typography} from '@material-ui/core'
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';

const Row = (props) => {
    const {row} = props
    const [open, setOpen] = useState(false)
    return (
        <React.Fragment>
            <TableRow>
                <TableCell>
                    {row.name}
                    <IconButton onClick={() => setOpen(!open)}>
                        {open ? <KeyboardArrowUpIcon/> : <KeyboardArrowDownIcon/> }
                    </IconButton>
                </TableCell>
                <TableCell>{row.cat1}</TableCell>
                <TableCell>{row.cat2}</TableCell>
                <TableCell>{row.cat3}</TableCell>
            </TableRow>

            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={4}>
                    <Collapse in={open}> 
                        <Box margin={2}>
                        <Typography variant="h6" gutterBottom>
                            More Data
                        </Typography>
                            <Table size="small">
                                <TableHead>
                                    <TableRow>
                                        <TableCell>test 1</TableCell>
                                        <TableCell>test 2</TableCell>
                                    </TableRow>
                                    
                                </TableHead>
                                <TableBody>
                                    <TableRow>
                                        <TableCell>test 1</TableCell>
                                        <TableCell>test 2</TableCell>
                                    </TableRow>
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
    // const {rows} = props
    // TEMP DATA
    let rows = [
        {id: 1, name: "resource 1", cat1: "test1", cat2: "test2",cat3: "test3"},
        {id: 2, name: "resource 2", cat1: "test1", cat2: "test2",cat3: "test3"},
        {id: 3, name: "resource 3", cat1: "test1", cat2: "test2",cat3: "test3"}
    ]

    return (
        <Paper elevation={3}>
            <Table>
                <TableHead>
                    <TableRow>
                        <TableCell>Resources </TableCell>
                        <TableCell>Category 1 </TableCell>
                        <TableCell>Category 2 </TableCell>
                        <TableCell>Category 3 </TableCell>

                    </TableRow>
                </TableHead>
                <TableBody>
                    {rows.map((row) => (
                        <Row key={row.id} row={row}/>
                    ))}
                </TableBody>
            </Table>
        </Paper>
        
    )
}

export default ResourcesTable