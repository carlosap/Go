import React, { createContext, useReducer, useEffect, useState } from 'react';
import { appReducer } from '../reducers/appReducer';
export const AppContext = createContext();

let socket = null;

const AppContextProvider = (props) => {

    const [isConnection, setConnection] = useState(false);

    const [state, dispatch] = useReducer(appReducer, {}, () => {
        const weather = localStorage.getItem('weather');
        const forecast = localStorage.getItem('forecast');
        const ipinfo = localStorage.getItem('ipinfo');
        const news = localStorage.getItem('news')
        
        return {
            weather: weather ? JSON.parse(weather) : {},
            forecast: forecast ? JSON.parse(forecast) : {},
            ipinfo: ipinfo ? JSON.parse(ipinfo) : {},
            news: news ? JSON.parse(news) : {}
        }
    })

    useEffect(() => {localStorage.setItem('weather', JSON.stringify(state.weather))}, [state.weather]);
    useEffect(() => {localStorage.setItem('forecast', JSON.stringify(state.forecast))}, [state.forecast]);
    useEffect(() => {localStorage.setItem('ipinfo', JSON.stringify(state.ipinfo))}, [state.ipinfo]);
    useEffect(() => {localStorage.setItem('news', JSON.stringify(state.news))}, [state.news]);
    
    const wsUpdateListener = () => {
        socket.onmessage = (e) => {
            console.log(e.data)
            try {
                const data = e.data ? JSON.parse(e.data) : {}
                setReducerState(data)
            } catch (error) {
                console.log(error)
            }

        };

    };


    const wsClosedListener = (e) => {

        // if (socket) {
        //     console.error(`Disconnected.....${isConnection}`);
        // }

        if (!isConnection) {
            socket = new WebSocket("ws://localhost:5000/ws");
            socket.onopen = () => {
                console.log("successfully connected");
                if (socket.readyState !== WebSocket.CLOSED || socket.readyState !== WebSocket.CONNECTING) {
                    //socket.send('hi from client');
                    setConnection(true)
                }

            };
            socket.onclose = (event) => {
                setConnection(false)
                console.log("socket closed connection: ", e);
                wsClosedListener();
            };

            window.setTimeout(wsUpdateListener);
        }
    }


    const setReducerState = (data) => {
        switch(data.msgtype) {
            case 'weather':
                dispatch({type: 'SET_WEATHER', payload: data});
                break;
            case 'forecast':
                dispatch({type: 'SET_FORECAST', payload: data});
                break;
            case 'ipinfo':
                dispatch({type: 'SET_IPINFO', payload: data});
                break;
            case 'news':
                dispatch({type: 'SET_NEWS', payload: data});
                break;

            default:
                break;
        }
    }
        
    useEffect(() => {wsClosedListener();})
    return(
        <AppContext.Provider value={{state, dispatch}}>
            {props.children}
        </AppContext.Provider>
    )
}

export default AppContextProvider;