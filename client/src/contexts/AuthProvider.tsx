import React, { useContext, useEffect, useState } from 'react';
import firebase from '../firebase';

type ContextProps = {
  user: firebase.User | null;
  authenticated: boolean;
  setUser: any;
  loadingAuthState: boolean;
};

const AuthContext = React.createContext<Partial<ContextProps>>({});

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }: any) => {
  const [user, setUser] = useState<ContextProps['user']>(null);
  const [loadingAuthState, setLoadingAuthState] = useState<
    ContextProps['loadingAuthState']
  >(true);

  useEffect(() => {
    firebase.auth().onAuthStateChanged((user: any) => {
      setUser(user);
      setLoadingAuthState(false);
    });
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        authenticated: user !== null,
        setUser,
        loadingAuthState,
      }}
    >
      {!loadingAuthState && children}
    </AuthContext.Provider>
  );
};
