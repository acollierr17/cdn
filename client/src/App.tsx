import React from 'react';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';
import { Helmet } from 'react-helmet';
import Dashboard from './containers/Dashboard';

export default function App() {
  return (
    <Router>
      <Helmet titleTemplate="Anthony • %s" defaultTitle="Anthony Collier">
        <html lang="en" />
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />

        {/* Primary Meta Tags  */}
        <meta
          name="description"
          content="My personal website and corner on the internet!"
        />
        <meta name="author" content="Anthony Collier" />
        <meta name="theme-color" content="#dd9323" />

        {/* Open Graph / FaceBook */}
        <meta property="og:type" content="website" />
        <meta property="og:url" content="https://acollier.dev/" />
        <meta
          property="og:description"
          content="My personal website and corner on the internet!"
        />
        <meta property="og:image" content="" />
      </Helmet>
      <Switch>
        <Route exact path="/" component={Dashboard} />
      </Switch>
    </Router>
  );
}
