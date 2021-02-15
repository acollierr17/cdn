import React from 'react';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';
import { Helmet } from 'react-helmet';
import Dashboard from './containers/Dashboard';
import Login from './containers/Login';
import NotFound from './components/NotFound';
import PrivateRoute from './components/PrivateRouter';
import { AuthProvider } from './contexts/AuthProvider';
import Profile from './containers/Profile';

export default function App() {
  return (
    <AuthProvider>
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
          <PrivateRoute exact path="/" component={Dashboard} />
          <PrivateRoute path="/profile" component={Profile} />
          <Route path="/login" component={Login} />
          <Route component={NotFound} />
        </Switch>
      </Router>
    </AuthProvider>
  );
}
