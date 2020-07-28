import React from 'react';
import News from '../components/News';

const NewsPage = () => {
    return (
        <div>
            <div className="content pure-u-1 pure-u-md-3-5">
                <News />
            </div>
            <div className="content pure-u-1 pure-u-md-2-5">
                <News />
            </div>
        </div>
    )
}

export default NewsPage;