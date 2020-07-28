import React, { useContext } from 'react';
import LeafletMap from '../components/LeafletMap';
import { AppContext } from '../contexts/AppContext';
import '../styles/news.css'

const IpInfo = () => {
    const { state } = useContext(AppContext)
    const ip = state.ipinfo;
    return (
        <div>
            <div className="posts">
                <div className="content pure-u-1 pure-u-md-3-5">
                    <h1 className="content-subhead">Continent</h1>
                    <section className="post">
                        <header className="post-header">
                            <img width="48" height="48" className="pure-img" src={ip.location.country_flag} />
                            <h2 className="post-title">{ip.continent_name} - {ip.country_name} ({ip.country_code})</h2>
                        </header>

                    </section>
                </div>
                <div className="content pure-u-1 pure-u-md-1-5">
                    <h1 className="content-subhead">Region</h1>
                    <section className="post">
                        <header className="post-header">
                            <h2 className="post-title">{ip.region_name} ({ip.region_code})</h2>
                        </header>

                    </section>
                </div>
                <div className="content pure-u-1 pure-u-md-1-5">
                    <h1 className="content-subhead">City</h1>
                    <section className="post">
                        <header className="post-header">
                            <h2 className="post-title">{ip.city}</h2>
                        </header>

                    </section>
                </div>
            </div>

            <div className="posts">
                <div className="content pure-u-1 pure-u-md-2-5">
                    <h1 className="content-subhead">Coordinates</h1>
                    <section className="post">
                        <header className="post-header">
                            <h2 className="post-title">Latitude:  {ip.latitude}</h2>
                        </header>
                        <header className="post-header">
                            <h2 className="post-title">Longitude:  {ip.longitude}</h2>
                        </header>
                    </section>
                </div>

                <div className="content pure-u-1 pure-u-md-3-5">
                    <h1 className="content-subhead">Map</h1>
                    <section className="post">
                    <div id="container">
                    MAP HERE...
                    </div>
                        
                    </section>
                </div>
                <div className="content pure-u-1 pure-u-md-1-5">
                    <h1 className="content-subhead">Language</h1>
                    <section className="post">
                        <header className="post-header">
                            <h2 className="post-title">{ip.location.languages[0].name}</h2>
                        </header>

                    </section>
                </div>
            </div>

        </div>


    )
}

export default IpInfo;