import firebase from 'firebase/app';
import 'firebase/firestore';
import 'firebase/auth';

import firebaseConfig from './config/firebase';

const app = firebase.initializeApp(firebaseConfig);

const firestore = app.firestore();
export const database = {
  tokens: firestore.collection('tokens'),
};
export const auth = app.auth();
