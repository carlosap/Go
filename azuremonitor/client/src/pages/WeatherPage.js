import React from 'react';
import Weather from '../components/Weather';
import Forecast from '../components/Forecast';

const WeatherPage = () => {

    return (
        <div>
            <h1>Weather</h1>
            <div>
                <Weather />
            </div>

            <h1>Forecast</h1>
            <div>
                <Forecast />
            </div>
        </div>
    )
}

export default WeatherPage;