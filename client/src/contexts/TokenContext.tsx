import React, { useReducer, useEffect } from 'react';
import { TokenActionTypes, tokenReducer } from '../reducers/token';
import { useAuth } from './AuthContext';
import { database } from '../firebase';

export const useToken = (uid: string = '') => {
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
