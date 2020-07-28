import React, { useContext, useState } from 'react';
import { AppContext } from '../contexts/AppContext';
import {Divider, Icon, Paper, Switch, Typography} from '@material-ui/core'
import LocationOnIcon from '@material-ui/icons/LocationOn';

const Weather = () => {
     const {state} = useContext(AppContext)
     console.log(state)
    // const currentWeather = weather.weather.weather[0] // Current weather within JSON

    // console.log(currentWeather)
    // const [isCelius, setIsCelcius] = useState(false)

    // // Convert from kelvin to desired unit
    // const getConvertedTemperature = (temp, trailingDecimals) => {
    //     // If fahrenheit
    //     if(!this.state.tempUnit){
    //         return ((temp - 273.15) * (9/5) + 32).toFixed(trailingDecimals)
    //     } else { // else return celcius 
    //         return (temp - 273.15).toFixed(trailingDecimals)
    //     }
    // }

    // // Get Cardinal Angle of Wind
    // const getCardinal = (angle) => {
    //     const degreePerDirection = 360 / 8;
      
    //     /** 
    //      * Offset the angle by half of the degrees per direction
    //      * Example: in 4 direction system North (320-45) becomes (0-90)
    //      */
    //     const offsetAngle = angle + degreePerDirection / 2;
      
    //     return (offsetAngle >= 0 * degreePerDirection && offsetAngle < 1 * degreePerDirection) ? "N"
    //       : (offsetAngle >= 1 * degreePerDirection && offsetAngle < 2 * degreePerDirection) ? "NE"
    //         : (offsetAngle >= 2 * degreePerDirection && offsetAngle < 3 * degreePerDirection) ? "E"
    //           : (offsetAngle >= 3 * degreePerDirection && offsetAngle < 4 * degreePerDirection) ? "SE"
    //             : (offsetAngle >= 4 * degreePerDirection && offsetAngle < 5 * degreePerDirection) ? "S"
    //               : (offsetAngle >= 5 * degreePerDirection && offsetAngle < 6 * degreePerDirection) ? "SW"
    //                 : (offsetAngle >= 6 * degreePerDirection && offsetAngle < 7 * degreePerDirection) ? "W"
    //                   : "NW";
    // }

    return (
        // <Paper className='WeatherCard-Container' elevation={3} square={false}>
        //     <div className="WeatherCard-Location">
        //         <LocationOnIcon style={{paddingTop:"5px", marginRight:'4px'}}/> 
        //         <Typography color='primary' variant='h4'>{this.props.city}</Typography>
        //     </div>
        //     <Divider variant="middle"/>

            
        //     <div className="WeatherCard-Description">
        //         <Typography variant='h5' color="textPrimary">
        //             {currentWeather.description ? 
        //             this.formatDescription(currentWeather.description) : null}
        //         </Typography>
        //     </div>

        //     <div className="WeatherCard-IconContainer">
        //         {this.getWeatherIcon()}
        //     </div>

        //     <div className='WeatherCard-Converter'>
        //         <Typography color="textPrimary">F°</Typography>
        //         <Switch checked={this.state.tempUnit} onChange={this.handleTempConversion}></Switch>
        //         <Typography color="textPrimary">C°</Typography>
        //     </div>

        //     <div className="WeatherCard-Temperature">
        //         <Typography variant='h4' color="textPrimary">{this.getConvertedTemperature(this.props.main.temp, 2)}</Typography>
        //         <Typography className="ml-4" variant='h4' color="textPrimary">°</Typography>
        //         <Typography variant='h4' color="textPrimary">
        //             {this.state.tempUnit ? 'C' : 'F'}
        //         </Typography>
        //     </div>

        //     <div className="WeatherCard-HighLow">
        //         <Typography variant='h6' color="textSecondary">High/Low - </Typography>
        //         <Typography style={{marginLeft: '5px'}} variant='h6' color="textSecondary">
        //             {this.getConvertedTemperature(this.props.main.temp_max, 0)}°/
        //             {this.getConvertedTemperature(this.props.main.temp_min, 0)}°
        //         </Typography> 
        //     </div>
            
        //     <Divider/>

        //     <div className="WeatherCard-Wind">
        //         <Icon style={{margin:"0 5px"}} className="meteocons" color="action">wind</Icon>
        //         <Typography style={{marginRight: '8px'}}color='textSecondary' variant='subtitle1'>
        //             {this.props.wind.speed} M/S
        //         </Typography>

        //         <Icon style={{margin:"0 5px"}} className="meteocons" color="action">compass</Icon>
        //         <Typography color='textSecondary' variant='subtitle1'>
        //             {this.getCardinal(this.props.wind.deg)}
        //         </Typography>
        //     </div>
        // </Paper>
        <div>
        {state.weather ? JSON.stringify(state.weather) : ''}
        </div>
    )
}

export default Weather;