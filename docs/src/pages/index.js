import React from 'react';
import clsx from 'clsx';
import Layout from '@theme/Layout';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import useBaseUrl from '@docusaurus/useBaseUrl';
import styles from './styles.module.css';

const features = [
    {
        title: 'Simple, Immutable Types',
        imageUrl: 'img/favicon.svg',
        description: (
            <>
                <p>A core premise of arr.ai is that the type system should be made as simple as
                    possible, but not simpler.</p>
                <p>Numbers, tuples and sets are sufficiently powerful to represent a rich and
                    diverse set of information structures, including arrays, string, dictionaries
                    and even functions.</p>
            </>
        ),
    },
    {
        title: 'Expressive Syntax',
        imageUrl: 'img/favicon.svg',
        description: (
            <>
                <p>Arr.ai is optimised for expressivity, <em>not</em> ease of learning. As a
                    functional language designed around relational algebra, it is unlike most other
                    languages. </p>

                <p>Investment in understanding fundamental concepts, operations and idioms is repaid
                    with a powerful tool that can (optimistically) reduce the volume of code you
                    need to write by an order of magnitude.</p>
            </>
        ),
    },
    {
        title: 'Hermetic runtime',
        imageUrl: 'img/favicon.svg',
        description: (
            <>
                Extend or customize your website layout by reusing React. Docusaurus can
                be extended while reusing the same header and footer.
            </>
        ),
    },
];

function Feature({imageUrl, title, description}) {
    const imgUrl = useBaseUrl(imageUrl);
    return (
        <div className={clsx('col col--4', styles.feature)}>
            {imgUrl && (
                <div className="text--center">
                    <img className={styles.featureImage} src={imgUrl} alt={title}/>
                </div>
            )}
            <h3>{title}</h3>
            <p>{description}</p>
        </div>
    );
}

function Home() {
    const context = useDocusaurusContext();
    const {siteConfig = {}} = context;
    return (
        <Layout
            title={`${siteConfig.title}`}
            description="The ultimate data engine">
            <header className={clsx('hero hero--primary', styles.heroBanner)}>
                <div className="container">
                    <h1 className="hero__title">{siteConfig.title}</h1>
                    <p className="hero__subtitle">{siteConfig.tagline}</p>
                    <div className={styles.buttons}>
                        <Link
                            className={clsx(
                                'button button--outline button--secondary button--lg',
                                styles.getStarted,
                            )}
                            to={useBaseUrl('docs/')}>
                            Get Started
                        </Link>
                    </div>
                </div>
            </header>
            <main>
                {features && features.length > 0 && (
                    <section className={styles.features}>
                        <div className="container">
                            <div className="row">
                                {features.map((props, idx) => (
                                    <Feature key={idx} {...props} />
                                ))}
                            </div>
                        </div>
                    </section>
                )}
            </main>
        </Layout>
    );
}

export default Home;
