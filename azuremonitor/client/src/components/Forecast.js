import React, { useContext } from 'react';
import { AppContext } from '../contexts/AppContext';

const Forecast = () => {
    const { state } = useContext(AppContext)
    return (
    <div>
        {state.forecast ? JSON.stringify(state.forecast) : ''}
    </div>
    )
}

export default Forecast;