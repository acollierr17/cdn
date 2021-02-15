import React, { useContext, useEffect, useState } from 'react';
import firebase from 'firebase/app';
import { auth } from '../firebase';

export type UserData = {
  email: string;
  password: string;
};

type ContextProps = {
  user: firebase.User | null;
  authenticated: boolean;
  setUser: React.Dispatch<React.SetStateAction<ContextProps['user']>>;
  loadingAuthState: boolean;
  login(data: UserData, history: Record<string, any>): Promise<void>;
  logout(history: Record<string, any>): Promise<void>;
};

const AuthContext = React.createContext<Partial<ContextProps>>({});

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }: any) => {
  const [user, setUser] = useState<ContextProps['user']>(null);
  const [loadingAuthState, setLoadingAuthState] = useState<
    ContextProps['loadingAuthState']
  >(true);

  const login = (data: UserData, history: Record<string, any>) =>
    auth.signInWithEmailAndPassword(data.email, data.password).then((res) => {
      setUser(res.user);
      history.push('/');
    });

  const logout = (history: Record<string, any>) =>
    auth.signOut().then(() => history.push('/login'));

  useEffect(() => {
    return firebase.auth().onAuthStateChanged((user: any) => {
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
        login,
        logout,
      }}
    >
      {!loadingAuthState && children}
    </AuthContext.Provider>
  );
};
