
import React from 'react';
import { Map, Marker, Popup, TileLayer } from 'react-leaflet'
import L from "leaflet"
//import '../styles/map.css'

const LeafletMap = ({ latitude, longitude, city, region }) => {
    const marker = new L.Icon({
        iconUrl: "https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-blue.png",
        iconSize: [15, 20],
        iconAnchor: [8, 18]
    })
    

    const getPosition = (lat, long) => {
        const latitude = lat ? lat : 0;
        const longitude = long ? long : 0;
        console.log(`${latitude} - ${longitude}`)
        return [0,0]
    } 
    
    return (
        <div>

                <Map center={[51.505,-0.09]} zoom={13}>
                    <TileLayer
                        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    />
                    <Marker
                        icon={marker}
                        position={[51.505,-0.09]}>
                        <Popup>
                            <div>
                                <div>Your Location:</div>
                                <div>{city}, {region}</div>
                            </div>
                        </Popup>
                    </Marker>
                </Map>
            </div>
       
    )
}

export default LeafletMap;