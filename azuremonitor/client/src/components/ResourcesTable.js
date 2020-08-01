import React from 'react'
import {makeStyles} from '@material-ui/core/styles'
import {Table, TableContainer, TableHead, TableBody, TableRow, TableCell} from '@material-ui/core'
import {Paper} from '@material-ui/core'

const ResourcesTable = (props) => {
    return (
        <TableContainer component={Paper}>
            <Table>
                <TableHead>
                    <TableRow>
                        <TableCell>Resources </TableCell>
                        <TableCell>Type </TableCell>
                        <TableCell>Some Category </TableCell>
                        <TableCell>Some Category </TableCell>

                    </TableRow>
                </TableHead>
            </Table>
        </TableContainer>
        
    )
}

export default ResourcesTable