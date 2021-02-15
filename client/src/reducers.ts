type ActionMap<M extends { [index: string]: any }> = {
  [Key in keyof M]: M[Key] extends undefined
    ? {
        type: Key;
      }
    : {
        type: Key;
        payload: M[Key];
      };
};

export enum TokenActionTypes {
  SET_TOKEN = 'SET_TOKEN',
}

type Token = {
  uid: string;
  token: string;
};

type TokenPayload = {
  [TokenActionTypes.SET_TOKEN]: Partial<Token>;
};

export type TokenActions = ActionMap<TokenPayload>[keyof ActionMap<TokenPayload>];

export const tokenReducer = (state: Token, { type, payload }: TokenActions) => {
  switch (type) {
    case 'SET_TOKEN':
      return {
        ...state,
        token: payload.token!,
      };
    default:
      return state;
  }
};
