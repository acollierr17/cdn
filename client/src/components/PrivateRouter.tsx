import React from 'react';
import { Redirect, Route } from 'react-router-dom';
import { useAuth } from '../contexts/AuthProvider';

export default function PrivateRoute({ component: Component, ...rest }: any) {
  const { authenticated, loadingAuthState } = useAuth();
  if (loadingAuthState) return <div>Loading...</div>;

  return (
    <Route
      {...rest}
      render={(props) =>
        authenticated ? (
          <Component {...props} />
        ) : (
          <Redirect
            to={{ pathname: '/login', state: { prevPath: rest.path } }}
          />
        )
      }
    />
  );
}
