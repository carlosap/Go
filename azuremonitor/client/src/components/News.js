import React, { useContext } from 'react';
import { AppContext } from '../contexts/AppContext';
import '../styles/news.css'

const News = () => {
    const { state } = useContext(AppContext)
    return (
        <div>
            <div className="posts">
                <h1 className="content-subhead">Pinned Post</h1>
                <section className="post">
                    <header className="post-header">
                        <img width="48" height="48" alt="Tilo Mitra&#x27;s avatar" className="post-avatar" src="/img/common/tilo-avatar.png" />
                        <h2 className="post-title">Introducing Pure</h2>
                        <p className="post-meta">
                            By Tilo Mitra
                        </p>
                    </header>
                    <div className="post-description">
                        <p>description</p>
                    </div>
                </section>
            </div>

        </div>

    )
}

export default News;