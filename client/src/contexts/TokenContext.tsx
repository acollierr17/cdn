import React, { createContext, useReducer, Dispatch, useEffect } from 'react';
import { TokenActionTypes, TokenActions, tokenReducer } from '../reducers';
import { useAuth } from './AuthProvider';
import { database } from '../firebase';

type InitialStateType = {
  uid: string;
  token: string;
};

const initialState = {
  uid: '',
  token: '',
};

const TokenContext = createContext<{
  state: InitialStateType;
  dispatch: Dispatch<TokenActions>;
}>({
  state: initialState,
  dispatch: () => null,
});

const TokenProvider: React.FC = ({ children }) => {
  const [state, dispatch] = useReducer(tokenReducer, initialState);

  return (
    <TokenContext.Provider value={{ state, dispatch }}>
      {children}
    </TokenContext.Provider>
  );
};

const useToken = (uid: string = '') => {
  const { user: currentUser } = useAuth();

  const [state, dispatch] = useReducer(tokenReducer, {
    uid,
    token: '',
  });

  useEffect((): any => {
    if (uid === null)
      return dispatch({
        type: TokenActionTypes.SET_TOKEN,
        payload: { token: '' },
      });

    return database.tokens
      .doc(currentUser!.uid)
      .get()
      .then((doc) => {
        dispatch({
          type: TokenActionTypes.SET_TOKEN,
          payload: { token: doc.data()!.token },
        });
      })
      .catch(() => {
        dispatch({
          type: TokenActionTypes.SET_TOKEN,
          payload: { token: '' },
        });
      });
  }, [currentUser]);

  useEffect(() => {
    return database.tokens
      .where('uid', '==', currentUser!.uid)
      .onSnapshot({ includeMetadataChanges: true }, (snapshot) => {
        if (snapshot.empty) return;
        const doc = snapshot.docs[0]?.data();
        dispatch({
          type: TokenActionTypes.SET_TOKEN,
          payload: { token: doc.token },
        });
      });
  }, [currentUser]);

  return state;
};

export { useToken, TokenProvider };
